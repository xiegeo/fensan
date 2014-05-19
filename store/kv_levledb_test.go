package store

import (
	"bytes"
	"os"
	"testing"
)

func TestKVL(t *testing.T) {
	path := ".TestKVL"
	kv, err := OpenLeveldb(path)
	if err != nil {
		t.Fatal(err)
	}
	testKL(kv, t)
	kv.Close()
	os.RemoveAll(path)
}

func testKL(kv KV, t *testing.T) {
	k := []byte{1}
	value := []byte{4, 5, 6}

	kv.Set(k, value)
	kv.Sync()
	v := kv.Get(k)
	if !bytes.Equal(value, v) {
		t.Error("can't get back set value")
	}
	kv.Delete(k)
	kv.Sync()
	v = kv.Get(k)
	if len(v) != 0 {
		t.Error("deleted value not removed")
	}
}
