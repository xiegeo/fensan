package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

//kvl is a KV backed by leveldb
type kvl struct {
	db *leveldb.DB
	ro *opt.ReadOptions
	b  *leveldb.Batch
}

func OpenLeveldb(path string) (KV, error) {
	o := &opt.Options{
		Compression: opt.NoCompression, //most data are hashes or small binary.
	}
	ro := &opt.ReadOptions{
		Strict: opt.StrictAll,
	}
	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		return nil, err
	}
	return &kvl{
		db: db,
		ro: ro,
		b:  new(leveldb.Batch),
	}, nil
}

func (kv *kvl) Get(key []byte) []byte {
	v, _ := kv.db.Get(key, kv.ro)
	//ignor ErrNotFound, len(v) should be 0 if so
	return v
}

func (kv *kvl) Set(key []byte, v []byte) {
	kv.b.Put(key, v)
}

func (kv *kvl) Delete(key []byte) {
	kv.b.Delete(key)
}

var kvlSync opt.WriteOptions = opt.WriteOptions{Sync: true}

func (kv *kvl) Sync() {
	err := kv.db.Write(kv.b, &kvlSync)
	if err != nil {
		panic(err)
	}
	kv.b.Reset()
}

func (kv *kvl) Close() error {
	kv.Sync()
	return kv.db.Close()
}

func (kv *kvl) GC(startAfterKey []byte, f func(key []byte, v []byte) (delete bool, stop bool)) {
	//todo
}
