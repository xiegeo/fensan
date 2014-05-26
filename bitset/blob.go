package bitset

import (
	"fmt"
	"os"
)

//Blob can be used to read/write parts of a long '[]byte' values
//without loading the full data.
//It can be backed by os.File, as it is a subset of os.File, except Size.
//All error reporting are turned off, as none are expected (panic otherwise).
type Blob interface {
	//length in bytes for this block
	Size() int64
	//ReadAt reads len(b) bytes from the File starting at byte offset off.
	ReadAt(b []byte, off int64)
	//WriteAt writes len(b) bytes to the File starting at byte offset off.
	//WriteAt, may or may not sync to disk, at any order.
	WriteAt(b []byte, off int64)
	//Sync commits the all WriteAt to stable storage.
	Sync()
	//Close closes the Block, allowing the freeing of resources. This Block can
	//not be used after Close.
	Close()
}

type fileBlob struct {
	f    *os.File
	size int64
}

func NewBlobFromFile(file *os.File, size int64) Blob {
	return &fileBlob{file, size}
}

func (f *fileBlob) Size() int64 {
	return f.size
}

func (f *fileBlob) ReadAt(b []byte, off int64) {
	assertInRange(b, off, f.size)
	n, err := f.f.ReadAt(b, off)
	if n != len(b) || err != nil {
		panic(fmt.Errorf("can't ReadAt:%v, %v, %v", off, n, err))
	}
}

func (f *fileBlob) WriteAt(b []byte, off int64) {
	assertInRange(b, off, f.size)
	n, err := f.f.WriteAt(b, off)
	if n != len(b) || err != nil {
		panic(fmt.Errorf("can't WriteAt:%v, %v, %v", off, n, err))
	}
}

func (f *fileBlob) Sync() {
	err := f.f.Sync()
	if err != nil {
		panic(err)
	}
}

func (f *fileBlob) Close() {
	err := f.f.Close()
	if err != nil {
		panic(err)
	}
}

func assertInRange(buf []byte, off int64, size int64) {
	if off < 0 {
		panic(fmt.Errorf("out of range:%v < 0", off))
	}
	if int64(len(buf))+off > size {
		panic(fmt.Errorf("out of range:%v + %v > %v", len(buf), off, size))
	}
}