//Store is a package for managing persistent local server data,
//including data blobs, metadata, user data, and server configerations.
//
//
//Large files are divided into blocks for internal server deduplication, this is
//invisble to the outside. Lower block size gives better deduplication, and
//higher block size gives faster sequential access and less metadata. (use 4MB for now)
//
//Store is designed to store files with the following possibilities:
//
//Importing a local file: data is read, hashed, and saved a block at a time. The
//hash and length is reported at the end.
//
//Requested to be saved locally: hash and length are given (along with locations
//to the server). Inner hash are downloaded to see what blocks are needed, then
//the blocks are downloaded.
//
//As such hashes maybe stored bottom up or top down, and when a tree is traversed,
//it may find missing nodes. An error condition that the server must handle or report
//to client.
package store

import (
	"io"

	ht "github.com/xiegeo/fensan/hashtree"
)

//KV is an interface for a []byte based key value store. For storing small pieces
//of metadata that are often updated.
type KV interface {
	//Get gets the value for the given key. It returns nil if KV does not contain the key.
	//
	//Warning: Get might not see what was Set or Deleted until after sync().
	//This is a temperary work around to easly warp some apis.
	Get(key []byte) []byte
	//Set sets the value for the given key. It overwrites any previous value for that key.
	//Do not set with len(v) == 0, use Delete instead.
	Set(key []byte, v []byte)
	//Delete deletes the value for the given key, it is a no-op if KV does not contain the key.
	Delete(key []byte)
	//Sync commits all the changes to stable storage. It blocks until done.
	Sync()
	//Close closes KV.
	Close() error
	GC(startAfterKey []byte, f func(key []byte, v []byte) (delete bool, stop bool))
}

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

//LV (Large Values) is an interface for a Block based data store. For
//storing large values with length known in advance and as part of key.
//The same key with two sizes represent two different Blobs.
type LV interface {
	//New creates a new block in LV, or nil if the key already exist.
	//size > 0
	New(key []byte, size int64) Blob
	//Get returns the Block referenced by key, or nil if key is not found.
	Get(key []byte, size int64) Blob
	//Move changes the key from from to to, and resizes to newSize.
	//If from does not exist, or to already exist, then an error is reported.
	//The Blob keyed by from should not be open.
	Move(oldKey []byte, oldSize int64, newKey []byte, newSize int64) error

	//Delete removes the Block repersented by key from LV.
	//Returns true if blob is not in LV after Delete.
	Delete(key []byte, size int64) (bool, error)

	Close() error
}

type FileState int

const (
	//FILE_NONE means we don't have this file
	FILE_NONE FileState = iota
	//FILE_PART means we have parts of this file
	FILE_PART
	//FILE_COMPLETE means we have all of this file
	FILE_COMPLETE
)

//MetaStore is a store for hash meta data on top of KV.
//Such as: Child hashes, data length/hash level, and TTL.
type MetaStore interface {
	//InnerHashMinLevel reports the amount of inner hash saved, everything at or
	//above the Level can be retrieved
	InnerHashMinLevel() ht.Level
	//GetInnerHashes reads inner hashes at level and offset, len(hs) should be a
	//multiple of hash length (32).
	//
	//An error is reported if parts of hs is unknown, parts known is still filled
	//with unknown parts as all zeros.
	//
	//If parts of hs are impossible by index range, it panics
	GetInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) error

	PutInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) (has ht.Nodes, complete bool, err error)

	//TTLGet gets the TTL of a file
	TTLGet(key HLKey) TTL
	//TTLSetAtleast updates the TTL to be at least coved to util, the total
	//increase is multiplied by key.Length() to return costs in storage time by
	//byteMonth.
	//Savings in deduplication can be refected in byteMonth.
	TTLSetAtleast(key HLKey, util TTL) (byteMonth int64)

	//Close closes MetaStore
	Close()
}

//Database is the full interface for all the data of a server,
//WIP
type Database interface {
	MetaStore
	//GetState checks if we have a file or not, or in progress
	GetState(key HLKey) FileState
	//GetAt reads len(b) bytes of file key from offset off.
	//returns error is b can't be read completely
	GetAt(key HLKey, b []byte, off int64) error
	//PutAt writes b to file from offset off.
	//
	//has returns the number of leaf nodes completed. This can report the progress.
	//Iff has is full, complete is true
	//
	//If hash checking fails, an err is reported. The server can us this information
	//to demote the source.
	PutAt(key HLKey, b []byte, off int64) (has ht.Nodes, complete bool, err error)

	//Import a file from reader
	ImportFromReader(r io.Reader) HLKey
}

type metaValue struct {
	ttl TTL
}
