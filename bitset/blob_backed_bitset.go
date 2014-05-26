package bitset

import (
	"fmt"
	"os"
)

const (
	blobWriteBlock     = 4096 //byte rage of changes packed into one blob write operation
	blobWriteBlockBits = blobWriteBlock * 8
)

/*
BlobBackedBitSet assumes that Blob (os.File) intelligently cache file reads,
and flushes *only* on demand. BitSets are nomally used as metadata, so
it should only flash after the main data flushes.



Currently, it appears that file read can be an expansive operation,
cached or not. Use MakeFullBuffered.
*/
type BlobBackedBitSet struct {
	blob    Blob
	bits    int
	changes map[int]map[int]bool //[block number][bit number in block] = set as
}

//OpenFileBacked is the same as NewFileBacked but accepts the fileName or path.
//It opens the old file, or create a new file
func OpenFileBacked(fileName string, capacity int) *BlobBackedBitSet {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return NewFileBacked(file, capacity)
}

/*
NewFileBacked returns a bitset backed by file,
If file size match capacity, use the file as bitmap.
Otherwise replace the file sized to capacity (in bits) and zero filled.
*/
func NewFileBacked(f *os.File, capacity int) *BlobBackedBitSet {
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	bytes := bytesForbits(capacity)
	fileSize := fi.Size()
	if fileSize > bytes {
		f.Truncate(0)
		fileSize = 0
	}
	if fileSize < bytes {
		f.Truncate(bytes)
		_, err := f.WriteAt(make([]byte, bytes), 0)
		if err != nil {
			panic(err)
		}
	}
	blob := NewBlobFromFile(f, bytes)
	return NewBlobBacked(blob, capacity)
}

func NewBlobBacked(b Blob, capacity int) *BlobBackedBitSet {
	return &BlobBackedBitSet{b, capacity, make(map[int]map[int]bool)}
}

func bytesForbits(bits int) int64 {
	return (int64(bits) + 7) / 8
}

func (b *BlobBackedBitSet) FileByteLength() int64 {
	return bytesForbits(b.bits)
}

func (b *BlobBackedBitSet) locateChange(key int) (bucket int, bit int) {
	checkIndex(key, b.bits)
	bucket = key / blobWriteBlockBits
	bit = key - (bucket * blobWriteBlockBits)
	return
}

func (b *BlobBackedBitSet) locateByteMask(key int) (bucket int, mask byte) {
	checkIndex(key, b.bits)
	bucket = key / 8
	mask = 1 << byte(key%8)
	return
}

func (b *BlobBackedBitSet) setAs(i int, v bool) {
	bucket, bit := b.locateChange(i)
	bmap, ok := b.changes[bucket]
	if !ok {
		bmap = make(map[int]bool)
		b.changes[bucket] = bmap
	}
	bmap[bit] = v
}

func (b *BlobBackedBitSet) Set(i int) {
	b.setAs(i, true)
}

func (b *BlobBackedBitSet) Unset(i int) {
	b.setAs(i, false)
}

func (b *BlobBackedBitSet) Get(i int) bool {
	bucket, bit := b.locateChange(i)
	v, ok := b.changes[bucket][bit]
	if ok {
		return v
	}
	oneByte := make([]byte, 1)
	maskBucket, mask := b.locateByteMask(i)
	b.blob.ReadAt(oneByte, int64(maskBucket))
	return (oneByte[0] & mask) != 0
}

func (b *BlobBackedBitSet) Capacity() int { return b.bits }

//Sync changes and Close the file, BitSet must not be used again
func (b *BlobBackedBitSet) Close() {
	b.Sync()
	b.blob.Close()
	//nil all internal pointers
	b.blob = nil
	b.changes = nil
	b.bits = -1
}

func (b *BlobBackedBitSet) HaveChanges() bool {
	return len(b.changes) > 0
}

//Sync first calls flush, then calls Sync on the underling blob.
func (b *BlobBackedBitSet) Sync() {
	b.Flush()
	b.blob.Sync()
}

//Flush Writes all changes to the blob
func (b *BlobBackedBitSet) Flush() {
	if !b.HaveChanges() {
		return
	}
	buffer := make([]byte, blobWriteBlock)
	buckets := (b.bits + blobWriteBlockBits - 1) / blobWriteBlockBits
	for i := 0; i < buckets; i++ {
		bmap := b.changes[i]
		delete(b.changes, i)
		if len(bmap) > 0 {
			starts := i * blobWriteBlock
			if i == buckets-1 {
				bufferSize := (b.bits+8-1)/8 - starts
				buffer = make([]byte, bufferSize)
			}
			b.blob.ReadAt(buffer, int64(starts))
			for k, v := range bmap {
				maskBucket, mask := b.locateByteMask(k)
				if v {
					buffer[maskBucket] |= mask
				} else {
					buffer[maskBucket] &^= mask
				}
			}
			b.blob.WriteAt(buffer, int64(starts))
		}
	}
}

func (b *BlobBackedBitSet) ReadAt(buf []byte, off int64) {
	if b.HaveChanges() {
		fmt.Println("warning: read will flush, please flush explicitly")
		b.Flush()
	}
	b.blob.ReadAt(buf, off)
}

func (b *BlobBackedBitSet) WriteAt(buf []byte, off int64) {
	if b.HaveChanges() {
		fmt.Println("warning: write will flushes, please flush explicitly")
		b.Flush()
	}
	b.blob.ReadAt(buf, off)
}

func (b *BlobBackedBitSet) ExportBytes() []byte {
	exp := make([]byte, b.FileByteLength())
	b.ReadAt(exp, 0)
	return exp
}
