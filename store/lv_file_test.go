package store

import (
	"bytes"
	"os"
	"testing"
)

func TestLVF(t *testing.T) {
	path := ".TestLVF"
	lv := OpenFolderLV(path)
	blob := lv.New([]byte{1}, 8)
	if blob != nil {
		blob.WriteAt([]byte{3, 4, 5}, 4)
		blob.Close()
	}

	blob = lv.Get([]byte{1}, 8)
	out := make([]byte, 8)
	blob.ReadAt(out, 0)
	if !bytes.Equal(out, []byte{0, 0, 0, 0, 3, 4, 5, 0}) {
		t.Error("unexpected:", out)
	}
	blob.Close()
	d, err := lv.Delete([]byte{1}, 8)
	if err != nil {
		t.Log(err)
	}
	if !d {
		t.Error("can't delete")
	}
	blob = lv.Get([]byte{1}, 8)
	if blob != nil {
		t.Error("deleted value not removed")
	}
	lv.Close()
	os.RemoveAll(path)
}
