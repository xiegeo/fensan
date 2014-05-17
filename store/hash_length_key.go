package store

import ht "github.com/xiegeo/fensan/hashtree"

//HLKey is the hash and length of some data. Used to look up data in a high level
//database.
//Even though length is unnecessary for uniqueness, it is an usefull meta to keep
//around without significant overhead, for code path conditions, sanity checking,
//and debuging.
//
//The length is always included in very link and network request with the hash,
//or for child nodes inferred from the size of parent and it's position, so we
//can always know the length when we know the hash, without using any lookup.
//
//Future: HLKeys can inclued hidden database internals for faster lookups when
//a HLKey is made by the database itself.
type HLKey interface {
	Hash() []byte
	Length() ht.Bytes
}

type hLKey struct {
	hash   []byte
	length ht.Bytes
}

func NewHLKey(hash []byte, length ht.Bytes) HLKey {
	return &hLKey{hash, length}
}

func (k *hLKey) Hash() []byte {
	return k.hash
}

func (k *hLKey) Length() ht.Bytes {
	return k.length
}
