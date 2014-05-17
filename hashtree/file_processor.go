package hashtree

import (
	"crypto/sha256"
	"hash"
)

const (
	//LeafBlockSize is the max size in bytes
	//of a data block on the leaf of a hash tree.
	LeafBlockSize = 1024
)

// fileDigest represents the partial evaluation of a file hash.
type fileDigest struct {
	len           int64            // processed length
	leaf          hash.Hash        // a hash, used for hashing leaf nodes
	leafBlockSize int64            // size of base block in bytes
	tree          CopyableHashTree // the digest used for inner and root nodes
}

type fileDigestSample struct {
	HashTree
}

//I is a sample of the standared HashTree, to make some methods accessable
//without creating a new HashTree. Do not use it's write or sum functions.
var I = fileDigestSample{NewFile()}

// Create the standard file tree hash using leaf blocks of LeafBlockSize (1kB)
// and "crypto/sha256", and inner hash using sha256 (244's IHV) without padding.
func NewFile() HashTree {
	return NewFile2(LeafBlockSize, sha256.New(), NewTree2(NoPad32bytes, ht_sha256block))
}

// Create any tree hash using leaf blocks of size and leaf hash,
// and inner hash using tree hash, the tree stucture is internal to the tree hash.
func NewFile2(leafBlockSize int64, leaf hash.Hash, tree CopyableHashTree) HashTree {
	d := new(fileDigest)
	d.len = 0
	d.leafBlockSize = leafBlockSize
	d.leaf = leaf
	d.tree = tree
	return d
}

func (d *fileDigest) Nodes(len int64) Nodes {
	if len == 0 {
		return 1
	}
	return Nodes((len-1)/d.leafBlockSize) + 1
}

func (d *fileDigest) SetInnerHashListener(l func(level Level, index Nodes, hash, left, right *H256)) {
	d.tree.SetInnerHashListener(l)
}

func (d *fileDigest) Size() int { return d.tree.Size() }

func (d *fileDigest) BlockSize() int        { return int(d.leafBlockSize) }
func (d *fileDigest) BlockSizeBytes() int64 { return d.leafBlockSize }

func (d *fileDigest) Reset() {
	d.tree.Reset()
	d.leaf.Reset()
	d.len = 0
}

func (d *fileDigestSample) Write(p []byte) (int, error) {
	panic("the sample can not be writen to")
}

func (d *fileDigest) Write(p []byte) (int, error) {
	startLength := int64(len(p))
	xn := d.len % d.leafBlockSize
	for int64(len(p))+xn >= d.leafBlockSize {
		writeLength := d.leafBlockSize - xn
		d.leaf.Write(p[0:writeLength])
		p = p[writeLength:]
		d.tree.Write(d.leaf.Sum(nil))
		d.leaf.Reset()
		xn = 0
	}
	if len(p) > 0 {
		d.leaf.Write(p)
	}
	d.len += startLength
	return int(startLength), nil
}

func (d *fileDigestSample) Sum(in []byte) []byte {
	panic("the sample can not be summed")
}

func (d *fileDigest) Sum(in []byte) []byte {
	if d.len%d.leafBlockSize != 0 || d.len == 0 {
		// Make a copy of d.tree so that caller can keep writing and summing.
		tree := d.tree.Copy()
		tree.Write(d.leaf.Sum(nil))
		return tree.Sum(in)
	}
	return d.tree.Sum(in)
}
