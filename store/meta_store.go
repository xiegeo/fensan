package store

import (
	"fmt"

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

func (m *metaStore) asserInRange(key HLKey, hs []byte, level ht.Level, off ht.Nodes) (n ht.Nodes, rebased ht.Level) {
	lhs := len(hs)
	n, r := ht.Nodes(lhs/hashSize), lhs%hashSize
	if n == 0 {
		panic("hs can't hold a hash")
	}
	if r != 0 {
		panic("hs is not multples of hashes")
	}

	lw := ht.LevelWidth(ht.I.Nodes(key.Length()), level)
	if off < 0 || off+n >= lw {
		panic(fmt.Errorf("offset out: %v < 0 || %v + %v >= %v", off, off, n, lw))
	}
	rebased = level - m.minLevel + 1
	if rebased < 1 {
		panic("can't request levels lower than min")
	}
	return
}

func (m *metaStore) getHashBlob(key HLKey) (Blob, ht.Nodes) {
	fileBlobs := ht.LevelWidth(ht.I.Nodes(key.Length()), m.minLevel-1)
	nodesInTree := ht.HashTreeSize(fileBlobs)
	return m.hashStore.Get(key.Hash(), hashSize*nodesInTree), fileBlobs
}

func (m *metaStore) GetInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) error {
	n, rebased := m.asserInRange(key, hs, level, off)
	blob, fileBlobs := m.getHashBlob(key)
	blob.ReadAt(hs, hashSize*ht.HashPosition(fileBlobs, rebased, off))
	for i := ht.Nodes(0); i < n; i++ { //check for unfilled hashes
		h := hs[i*hashSize : hashSize]
		for j := 0; j < hashSize; j++ {
			if h[j] != 0 {
				goto next
			}
		}
		return fmt.Errorf("hash incomplete")
	next:
	}
	return nil
}

func (m *metaStore) PutInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) (has ht.Nodes, complete bool, err error) {
	panic("unimplemented")
	/*
		n, rebased := m.asserInRange(key, hs, level, off)
		blob, fileBlobs := m.getHashBlob(key)
		if n != fileBlob {
			return 0, false, fmt.Errorf("only full hash puts are supported for now. TODO: use bitset to support it")
		}
		//todo: check and fill
	*/
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
	return int64(buyMonth) * key.Length()
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
