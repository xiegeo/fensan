package store

import (
	"code.google.com/p/gogoprotobuf/proto"
	ht "github.com/xiegeo/fensan/hashtree"
	"github.com/xiegeo/fensan/pb"
)

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
	//Hash returns the hash of refereced. Do not modify the returned contests
	GetHash() []byte
	//Length returns the length of refereced file in bytes
	GetLength() int64
	//FullBytes is used when a key need to include length, as an attacker might
	//claim the existence of a file of the same hash but different size.
	//Do not modify the returned contests
	FullBytes() []byte
}

type hLKey struct {
	fullBytes []byte
	length    int64
}

//NewHLKey create a new HLKey, the hash is deep copied
func NewHLKey(hash []byte, length int64) HLKey {
	fb := make([]byte, ht.HashSize+8)
	copy(fb[:ht.HashSize], hash)
	littleEndianPutUint64(fb[ht.HashSize:], uint64(length))
	return &hLKey{fb, length}
}

func (k *hLKey) Proto() proto.Message {
	return pb.NewStaticIdFromFace(k)
}

func (k *hLKey) FullBytes() []byte {
	return k.fullBytes
}

func (k *hLKey) GetHash() []byte {
	return k.fullBytes[:ht.HashSize]
}

func (k *hLKey) GetLength() int64 {
	return k.length
}

//from encoding/binary/binary.go func (littleEndian) PutUint64
func littleEndianPutUint64(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

//from encoding/binary/binary.go func (littleEndian) Uint64
/* unused
func littleEndianUint64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}
*/
