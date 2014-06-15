package store

import (
	"bytes"
	"fmt"

	"github.com/xiegeo/fensan/bitset"
	ht "github.com/xiegeo/fensan/hashtree"
)

type metaStore struct {
	minLevel  ht.Level
	ttlStore  KV
	hashStore LV
}

const BlobSize = 4 << 20 //4MByte blocks
const hashSize = ht.HashSize

func OpenMetaStore(path string) (MetaStore, error) {
	minLevel := ht.Levels(BlobSize / ht.LeafBlockSize) // level 13
	ttlStore, err := OpenLeveldb(path + "/m_ttl")
	if err != nil {
		return nil, err
	}
	hashStore := OpenFolderLV(path + "/m_hash")
	return &metaStore{minLevel, ttlStore, hashStore}, nil
}

func (m *metaStore) InnerHashMinLevel() ht.Level {
	return m.minLevel
}

func (m *metaStore) asserInRange(key HLKey, hs []byte, level ht.Level, off ht.Nodes) (n, lw ht.Nodes, rebased ht.Level) {
	lhs := len(hs)
	n, r := ht.Nodes(lhs/hashSize), lhs%hashSize
	if n == 0 {
		panic("hs can't hold a hash")
	}
	if r != 0 {
		panic("hs is not multples of hashes")
	}

	lw = ht.LevelWidth(ht.I.Nodes(key.GetLength()), level)
	if off < 0 || off+n >= lw {
		panic(fmt.Errorf("offset out: %v < 0 || %v + %v >= %v", off, off, n, lw))
	}
	rebased = level - m.minLevel + 1
	if rebased < 1 {
		panic("can't request levels lower than min")
	}
	return
}

func (m *metaStore) mixedBlobSizes(key HLKey) (fileBlobs ht.Nodes, blobBytes, hashBytes, treeSize int64) {
	fileBlobs = ht.LevelWidth(ht.I.Nodes(key.GetLength()), m.minLevel-1)
	hashBytes = int64(fileBlobs) * hashSize
	treeSize = ht.HashTreeSize(fileBlobs)
	blobBytes = hashBytes + (treeSize+7)/8 + bitset.CountBytes
	return
}

func (m *metaStore) getHashBlob(key HLKey) (ht.Nodes, Blob, *bitset.CountingBitSet, bitset.Closer) {
	fileBlobs, blobBytes, hashBytes, treeSize := m.mixedBlobSizes(key)
	mixed := m.hashStore.Get(key.GetHash(), blobBytes)
	if mixed == nil {
		mixed = m.hashStore.New(key.GetHash(), blobBytes)
	}
	hashes, countingBlob := bitset.SplitBlob(mixed, hashBytes)
	countingBlob = bitset.MakeFullBuffered(countingBlob)
	counting := bitset.NewCounting(countingBlob, int(treeSize))
	return fileBlobs, hashes, counting, mixed
}

func (m *metaStore) GetInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) error {
	_, _, rebased := m.asserInRange(key, hs, level, off)
	fileBlobs, hashes, countingSet, closer := m.getHashBlob(key)
	defer closer.Close()
	if countingSet.Count() != int(fileBlobs) {
		return fmt.Errorf("hash incomplete")
	}
	hashes.ReadAt(hs, hashSize*ht.HashPosition(fileBlobs, rebased, off))
	return nil
}

func (m *metaStore) PutInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) (has ht.Nodes, complete bool, err error) {
	_, lw, rebased := m.asserInRange(key, hs, level, off)
	fileBlobs, hashes, countingSet, closer := m.getHashBlob(key)
	defer closer.Close()
	if countingSet.Count() == int(fileBlobs) {
		//already done, so it's a no-op
		return ht.Nodes(countingSet.Count()), true, nil
	}

	writeHash := func(rebasedL ht.Level, woff ht.Nodes, hash, left, right *ht.H256) {
		n := int(ht.HashNumber(fileBlobs, rebasedL, woff))
		woffBytes := int64(n * hashSize)
		hashes.WriteAt(hash.ToBytes(), woffBytes)
		countingSet.Set(n)
	}

	rootBuffer := make([]byte, hashSize)
	splited := ht.SplitLocalSummable(hs, hashSize, lw, off)
	c := ht.NewNoPadTree()
	var sum []byte
	for i := 0; i < len(splited); i++ {
		s := splited[i]
		hashHeight := ht.Levels(ht.Nodes(len(s) / hashSize))
		rebasedRootLevel := rebased + hashHeight
		rootOff := off >> uint(hashHeight)
		rootPosition := int(ht.HashNumber(fileBlobs, rebasedRootLevel, rootOff))
		if rootPosition == countingSet.Capacity() {
			copy(rootBuffer, key.GetHash()) //root is the key
		} else if !countingSet.Get(rootPosition) {
			goto next //don't have the root, skiped
		} else {
			hashes.ReadAt(rootBuffer, int64(rootPosition)*hashSize)
		}
		c.Write(s)
		sum = c.Sum(nil)
		c.Reset()
		if bytes.Equal(sum, rootBuffer) {
			//hashes verified, good for saving
			c.SetInnerHashListener(func(l ht.Level, hoff ht.Nodes, hash, left, right *ht.H256) {
				rebasedL := l + rebased
				woff := hoff + rootOff<<uint32(rebasedRootLevel-l)
				if rebasedL == rebasedRootLevel {
					return
				}
				writeHash(rebasedL, woff, hash, left, right)
				//propagate down nodes with single branch
				for rebasedL > 0 &&
					woff+1 == ht.LevelWidth(fileBlobs, rebasedL) &&
					woff%2 == 0 {
					rebasedL--
					woff = ht.LevelWidth(fileBlobs, rebasedL) - 1
					writeHash(rebasedL, woff, hash, left, right)
				}
			})
			c.Write(s)
			c.Sum(nil)
			c.SetInnerHashListener(nil)
			c.Reset()
		}
	next:
		off += ht.Nodes(len(s) / hashSize)
	}
	hashes.Sync()
	countingSet.Sync()
	return ht.Nodes(countingSet.Count()), countingSet.Full(), nil
}

func (m *metaStore) TTLGet(key HLKey) TTL {
	v := m.ttlStore.Get(key.FullBytes())
	if v == nil {
		return TTLLongAgo
	}
	return TTLFromBytes(v)
}
func (m *metaStore) TTLSetAtleast(key HLKey, freeFrom, until TTL) (byteMonth int64) {
	old := m.TTLGet(key)
	if old >= until {
		return 0
	}
	if freeFrom < old {
		freeFrom = old
	}
	buyMonth := freeFrom.MonthUntil(until)
	m.ttlStore.Set(key.FullBytes(), until.Bytes())
	return int64(buyMonth) * key.GetLength()
	//todo: check dedup
}

func (m *metaStore) Close() error {
	err := m.ttlStore.Close()
	err2 := m.hashStore.Close()
	if err != nil || err2 != nil {
		return fmt.Errorf("fail close meta store, ttlStore err: %v; hashStore err: %v", err, err2)
	}
	return nil
}
