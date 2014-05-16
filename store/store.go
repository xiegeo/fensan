//store is a package for managing persistent local server data,
//including data blobs, metadata, user data, and server configerations.
package store

import ht "github.com/xiegeo/fensan/hashtree"

//KV is an interface for a []byte based key value store.
type KV interface {
	Get(key []byte) []byte
	Set(key []byte, v []byte)
	Delete(key []byte)
	GC(startAfterKey []byte, f func(key []byte, v []byte) (delete bool, stop bool))
}

//MetaStore is a store for hash meta data on top of KV.
//Such as: Child hashes, data length/hash level, and TTL.
//MetaStore has large number of hashes as keys for look ups and deduplication.
type MetaStore interface {
}

type metaValue struct {
	ttl    TTL
	length ht.Bytes
	left   []byte
	right  []byte
}
