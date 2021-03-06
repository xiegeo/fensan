package store

import (
	"fmt"
	"os"
)

type folderLV struct {
	root       string
	permission os.FileMode
}

func OpenFolderLV(root string) LV {
	return &folderLV{root, 0777}
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
			return NewBlobFromFile(newFile, size)
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
	return NewBlobFromFile(opened, size)
}

func (f *folderLV) Move(oldKey []byte, oldSize int64, newKey []byte, newSize int64) error {
	panic("not implemented")
}

func (f *folderLV) Delete(key []byte, size int64) (bool, error) {
	_, file := byteToFile(f.root, key, size)
	err := os.Remove(file)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return true, err
	}
	return false, err
}

func (f *folderLV) Close() error {
	return nil //no-op
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
