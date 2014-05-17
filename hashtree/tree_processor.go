package hashtree

import (
	"fmt"
	"hash"
	"io"
)

//sha256 is 32 bytes
const HashSize = 32

//the max depth of a hash tree
const MaxLevel = 64

type H256 [8]uint32 //the internal hash

//FromBytes create a hash from bytes
//bytes must have a length of HashSize (32)
func FromBytes(bytes []byte) *H256 {
	var h H256
	for i := 0; i < 8; i++ {
		j := i * 4
		h[i] = uint32(bytes[j])<<24 | uint32(bytes[j+1])<<16 | uint32(bytes[j+2])<<8 | uint32(bytes[j+3])
	}
	return &h
}

//ToBytes reads out bytes from a hash
//bytes will have a length of HashSize (32)
func (h *H256) ToBytes() []byte {
	bytes := make([]byte, 32)
	for i, s := range h {
		bytes[i*4] = byte(s >> 24)
		bytes[i*4+1] = byte(s >> 16)
		bytes[i*4+2] = byte(s >> 8)
		bytes[i*4+3] = byte(s)
	}
	return bytes
}

type HashTree interface {
	hash.Hash

	//Nodes returns the number of nodes on the
	//bottom level of a hash tree covering len
	//bytes of data.
	//The padder is used by this function.
	Nodes(len int64) Nodes

	//SetInnerHashListener set a listener that receive
	//callbacks everytime an inner hash is calculated.
	//
	//On level 1, left and right child hashes are nil.
	//
	//On right-most nodes that are promoted without hashing,
	//left child is the same hash, and right is nil.
	SetInnerHashListener(l func(level Level, index Nodes, hash, left, right *H256))
}

//CopyableHashTree allows copying of the internal state.
//Because the Sum function might pad, but can't reset.
type CopyableHashTree interface {
	HashTree
	Copy() CopyableHashTree
}

// treeDigest represents the partial evaluation of a hashtree.
type treeDigest struct {
	x                 [HashSize]byte               // unprocessed bytes
	xn                int                          // length of x
	len               int64                        // processed length
	stack             [MaxLevel]*H256              // partial hashtree of more height then ever needed
	sn                Level                        // top of stack, depth of tree
	padder            func(d io.Writer, len int64) // the padding function
	compressor        func(l, r *H256) *H256       // 512 to 256 hash function
	innerHashListener func(level Level, index Nodes, hash, left, right *H256)
	innersCounter     [MaxLevel]Nodes
}

func NewTree() CopyableHashTree {
	return NewTree2(ZeroPad32bytes, ht_sha256block)
}

func NewNoPadTree() CopyableHashTree {
	return NewTree2(NoPad32bytes, ht_sha256block)
}

// Create a binary tree hash using padder and compressor.
// Padder mush pad to intervals of 256 bits.
// Compressor mush hash 2 H256s to 1.
func NewTree2(padder func(d io.Writer, len int64), compressor func(l, r *H256) *H256) CopyableHashTree {
	d := new(treeDigest)
	d.Reset()
	d.padder = padder
	d.compressor = compressor
	return d
}
func (d *treeDigest) Copy() CopyableHashTree {
	d0 := *d
	return &d0
}

type adder struct{ int64 }

//Write for adder only increment adder by length of input
func (a *adder) Write(p []byte) (length int, nil error) {
	length = len(p)
	a.int64 += int64(length)
	return
}

func (d *treeDigest) Nodes(len int64) Nodes {
	a := &adder{len}
	d.padder(a, len)
	return Nodes(a.int64)
}

func (d *treeDigest) SetInnerHashListener(l func(level Level, index Nodes, hash, left, right *H256)) {
	d.innerHashListener = l
}

func (d *treeDigest) Size() int { return HashSize }

func (d *treeDigest) BlockSize() int { return HashSize }

func (d *treeDigest) Reset() {
	d.xn = 0
	d.len = 0
	d.stack = [64]*H256{nil}
}
func (d *treeDigest) Write(p []byte) (startLength int, err error) {
	startLength = len(p)
	for len(p)+d.xn >= HashSize {
		for i := 0; i < HashSize-d.xn; i++ {
			d.x[d.xn+i] = p[i]
		}
		p = p[HashSize-d.xn:]
		d.xn = 0
		d.writeStack(0, FromBytes(d.x[:]), nil, nil)
	}
	if len(p) > 0 {
		for i := 0; i < len(p); i++ {
			d.x[d.xn+i] = p[i]
		}
		d.xn += len(p)
	}
	d.len += int64(startLength)
	return startLength, nil
}

func (d *treeDigest) listenInner(l Level, h, left, right *H256) {
	if d.innerHashListener != nil {
		d.innerHashListener(l, d.innersCounter[l], h, left, right)
	}
	d.innersCounter[l]++
}

func (d *treeDigest) writeStack(level Level, node, l, r *H256) {
	d.listenInner(level, node, l, r)
	if d.sn == level {
		d.stack[level] = node
		d.sn++
	} else if d.stack[level] == nil {
		d.stack[level] = node
	} else {
		last := d.stack[level]
		d.stack[level] = nil
		d.writeStack(level+1, d.compressor(last, node), last, node)
	}
}

func (d0 *treeDigest) Sum(in []byte) []byte {
	// Make a copy of d0 so that caller can keep writing and summing.
	d := *d0
	d.padder(&d, d.len)

	if d.xn != 0 {
		panic(fmt.Sprintf("d.xn = %d", d.xn))
	}

	var right *H256
	i := Level(0)
	for ; right == nil; i++ {
		right = d.stack[i]
	}
	if i < d.sn {
		d.listenInner(i, right, right, nil)
	}
	for ; i < d.sn; i++ {
		left := d.stack[i]
		if left != nil {
			oldR := right
			right = d.compressor(left, right)
			d.listenInner(i+1, right, left, oldR)
		} else {
			d.listenInner(i+1, right, right, nil)
		}
	}

	return append(in, right.ToBytes()...)
}

//ZeroPad32bytes pads with 0 or more of bytes 0x00.
//Use this when the size is known externally and
//the shortness of content makes padding the size
//too expansive.
func ZeroPad32bytes(d io.Writer, len int64) {
	padSize := (32 - (len % 32)) % 32
	if len == 0 {
		padSize = 32
	}
	d.Write(make([]byte, padSize))
}

//NoPad32bytes is a special case of ZeroPad32bytes
//when content is known to be in one or more multiples
//of 32 byte blocks, panic otherwise
func NoPad32bytes(d io.Writer, len int64) {
	if len%32 != 0 || len == 0 {
		panic(fmt.Sprintf("need padding of %v bytes for length of %v", 32-len%32, len))
	}
}
