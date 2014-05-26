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
	f            *os.File
	size         int64
	suppressSync bool
}

func NewBlobFromFile(file *os.File, size int64) Blob {
	return &fileBlob{file, size, true}
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
	f.suppressSync = false
}

func (f *fileBlob) Sync() {
	if f.suppressSync {
		return
	}
	f.suppressSync = true
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
	f.f = nil
}

func assertInRange(buf []byte, off int64, size int64) {
	if off < 0 {
		panic(fmt.Errorf("out of range:%v < 0", off))
	}
	if int64(len(buf))+off > size {
		panic(fmt.Errorf("out of range:%v + %v > %v", len(buf), off, size))
	}
}

type subBlob struct {
	blob  Blob
	start int64
	size  int64
}

func SplitBlob(b Blob, at int64) (left, right Blob) {
	if at <= 0 || at >= b.Size() {
		panic("index out of range")
	}
	left = &subBlob{b, 0, at}
	right = &subBlob{b, at, b.Size() - at}
	return
}

func (s *subBlob) Size() int64 {
	return s.size
}

func (s *subBlob) ReadAt(b []byte, off int64) {
	assertInRange(b, off, s.size)
	s.blob.ReadAt(b, off+s.start)
}

func (s *subBlob) WriteAt(b []byte, off int64) {
	assertInRange(b, off, s.size)
	s.blob.WriteAt(b, off+s.start)
}

func (s *subBlob) Sync() {
	s.blob.Sync()
}

func (s *subBlob) Close() {
	s.blob = nil
}

type fullBufferBlob struct {
	blob Blob
	buf  []byte
}

func MakeFullBuffered(blob Blob) Blob {
	buf := make([]byte, blob.Size())
	blob.ReadAt(buf, 0)
	return &fullBufferBlob{blob, buf}
}

func (f *fullBufferBlob) Size() int64 {
	return f.blob.Size()
}

func (f *fullBufferBlob) ReadAt(b []byte, off int64) {
	assertInRange(b, off, f.Size())
	copy(b, f.buf[off:])
}

func (f *fullBufferBlob) WriteAt(b []byte, off int64) {
	f.blob.WriteAt(b, off)
	copy(f.buf[off:], b)
}

func (f *fullBufferBlob) Sync() {
	f.blob.Sync()
}

func (f *fullBufferBlob) Close() {
	f.blob.Close()
	f.blob = nil
	f.buf = nil
}
