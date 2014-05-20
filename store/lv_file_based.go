package store

import (
	"fmt"
	"os"
)

type fileBlob struct {
	f    *os.File
	size int64
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

type folderLV struct {
	root       string
	permission os.FileMode
}

func (f *folderLV) New(key []byte, size int64) Blob {
	folder, file := byteToFile(f.root, key, size)
	os.MkdirAll(folder, f.permission)
	fi, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			newFile, err := os.Create(file)
			if err != nil {
				panic(err)
			}
			newFile.Truncate(size)
			return &fileBlob{newFile, size}
		} else {
			panic(err)
		}
	} else if fi.Size() == 0 {
		//if it crashed last time between Create and Truncate
		os.Truncate(file, size)
		return f.get(file, size)
	} else if fi.Size() == size {
		return nil
	} else {
		panic(fmt.Errorf("%v != %v", fi.Size(), size))
	}
}

func (f *folderLV) Get(key []byte, size int64) Blob {
	_, file := byteToFile(f.root, key, size)
	return f.get(file, size)
}

func (f *folderLV) get(file string, size int64) Blob {
	opened, err := os.OpenFile(file, os.O_RDWR, f.permission)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		panic(err)
	}
	fi, _ := opened.Stat()
	diskSize := fi.Size()
	if diskSize == 0 {
		return nil
	}
	if diskSize != size {
		panic("wrong size")
	}
	return &fileBlob{opened, size}
}

func (f *folderLV) Move(oldKey []byte, oldSize int64, newKey []byte, newSize int64) error {
	panic("not implemented")
}

func (f *folderLV) Delete(key []byte, size int64) {
	_, file := byteToFile(f.root, key, size)
	os.Remove(file)
}

const folderKeySize = 2

func byteToFile(root string, key []byte, size int64) (folder string, file string) {
	ks := make([]byte, folderKeySize) //ks is something from key in folderKeySize bytes
	if len(key) < folderKeySize {
		//TODO:key to short, need to extend, just use all zeros for now
	} else {
		ks = key[:folderKeySize]
		//TODO: "hash" the full key with a salt, so clients can't make just keys of same bucket.
	}

	folder = fmt.Sprintf("%s/%x/%x", root, ks[0], ks[1])
	file = fmt.Sprintf("%s/%x-%x", folder, size, key)
	return
}
