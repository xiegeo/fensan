package store

//HLKey is the hash and length of some data. Used to look up data in a high level
//database.
//Even though length is unnecessary for uniqueness, it is an usefull meta to keep
//around without significant overhead, for code path conditions, sanity checking,
//and debuging.
//
//Before a file can be fully checked, the length can not be conformed, it is to possible
//to have two downloads in progress, with the same hash but different length, where
//at least one can never be competed. So for all network communications, the length
//is nessary to go with the hash.
//
//The length is always included in very link and network request with the hash,
//or for child nodes inferred from the size of parent and it's position, so we
//can always know the length when we know the hash, without using any additional lookup.
//
type HLKey interface {
	Hash() []byte
	Length() int64
}

type hLKey struct {
	hash   []byte
	length int64
}

func NewHLKey(hash []byte, length int64) HLKey {
	return &hLKey{hash, length}
}

func (k *hLKey) Hash() []byte {
	return k.hash
}

func (k *hLKey) Length() int64 {
	return k.length
}
