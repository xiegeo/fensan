package bitset

import (
	"bytes"
	"encoding/binary"
	"log"
)

const (
	countBits  = 64
	countBytes = countBits / 8
)

type CountingBitSet struct {
	BlobBackedBitSet
	count int64
}

func (c *CountingBitSet) readCount() int64 {
	count := int64(0)
	buf := make([]byte, countBytes)
	c.blob.ReadAt(buf, int64(c.Capacity()+7)/8)
	err := binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &count)
	if err != nil {
		panic(err)
	}
	return count
}

func (c *CountingBitSet) writeCount(count int64) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, count)
	if err != nil {
		panic(err)
	}
	c.blob.WriteAt(buf.Bytes(), int64(c.Capacity()+7)/8)
}

func (c *CountingBitSet) recount() int64 {
	log.Println("recount started on bitset of capacity:", c.Capacity(),
		" this should only happen after a crash")
	if len(c.changes) != 0 {
		panic("recount can not be used with changes")
	}
	count := int64(0)
	buffer := make([]byte, blobWriteBlock)
	buckets := (c.Capacity() + blobWriteBlockBits - 1) / blobWriteBlockBits
	for i := 0; i < buckets; i++ {
		starts := i * blobWriteBlock
		if i == buckets-1 {
			bufferSize := (c.Capacity()+8-1)/8 - starts
			buffer = make([]byte, bufferSize)
		}
		c.blob.ReadAt(buffer, int64(starts))
		for _, b := range buffer {
			for i := uint(0); i < 8; i++ {
				if b&(1<<i) != 0 {
					count++
				}
			}
		}
	}
	return count
}

func OpenCountingFileBacked(fileName string, capacity int) *CountingBitSet {
	fileBacked := OpenFileBacked(fileName, capacity+countBits)
	return NewCounting(fileBacked.blob, capacity)
}

func NewCounting(blob Blob, capacity int) *CountingBitSet {
	if blob.Size() != int64(capacity+countBits+7)/8 {
		panic("blob of wrong size")
	}
	counting := &CountingBitSet{*NewBlobBacked(blob, capacity+countBits), 0}
	count := counting.readCount()
	if count == -1 {
		count = counting.recount()
	}
	counting.count = count
	return counting
}

func (c *CountingBitSet) Capacity() int { return c.bits - countBits }

func (c *CountingBitSet) Count() int {
	return int(c.count)
}

func (c *CountingBitSet) Full() bool {
	return c.Count() == c.Capacity()
}

func (c *CountingBitSet) Set(i int) {
	if c.Get(i) {
		return
	} else {
		c.count++
		c.BlobBackedBitSet.Set(i)
		if CHECK_INTEX && int(c.count) > c.Capacity() {
			panic("count > c.Capacity()")
		}
	}
}

func (c *CountingBitSet) Unset(i int) {
	if !c.Get(i) {
		return
	} else {
		c.count--
		c.BlobBackedBitSet.Unset(i)
		if CHECK_INTEX && c.count < 0 {
			panic("count < 0")
		}
	}
}

func (c *CountingBitSet) Sync() {
	c.Flush()
	c.BlobBackedBitSet.Sync()
	c.writeCount(c.count)
}

func (c *CountingBitSet) Flush() {
	if len(c.changes) == 0 {
		return
	}
	c.writeCount(-1) // mark count as dirty
	c.BlobBackedBitSet.Flush()
	//Flush does not guarantee write to disk, only write back count after sync
}

func (c *CountingBitSet) Close() {
	c.Sync()
	c.BlobBackedBitSet.Close()
}

func (c *CountingBitSet) DataByteLength() int64 {
	return int64(c.Capacity()+8-1) / 8
}

func (c *CountingBitSet) ExportBytes() []byte {
	return c.BlobBackedBitSet.ExportBytes()[:c.DataByteLength()]
}

func (c *CountingBitSet) ToSimple() *SimpleBitSet {
	return NewSimpleFromBytes(c.Capacity(), c.ExportBytes())
}
