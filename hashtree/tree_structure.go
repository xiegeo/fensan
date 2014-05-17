package hashtree

import (
	"math"
)

//Level is the depth of tree counted from bottom up.
//The lowest level is 0.
//Level is signed to allow representation of deltas.
type Level int

//int64 is the size in bytes.
//It is signed to allow representation of deltas.
//It is int64 as in file io api's usage.
//Since Level and Nodes are typed, bytes using int64 directly does not confuse
//with Level and Nodes.
//type int64 int64

//Nodes is the number of Nodes in a level,
//or the index of a node in a level.
//Nodes is signed to allow representation of deltas.
type Nodes int

//highBitLoc returns the location of the heighest bit, 1 based, so that 0 is all zeros
func highBitLoc(n uint32) uint32 {
	return uint32(math.Ilogb(float64(n)*2 + 1))
}

//highBitMask returns the largest power of 2 less or equal n
func highBitMask(n uint32) uint32 {
	return 1 << (highBitLoc(n) - 1) // 1 << -1 (or max uint32) = 0
}

func highBitMask64(n uint64) uint64 {
	if n >= 1<<32 {
		return uint64(highBitMask(uint32(n>>32)) << 32)
	} else {
		return uint64(highBitMask(uint32(n)))
	}

}

//Levels return the number of Levels (1 or more) of a hash
//tree with n leaf nodes
func Levels(n Nodes) Level {
	return Level(highBitLoc(uint32(n-1)) + 1)
}

//LevelWidth returns the number of nodes in a level,
//with the requested number of level (0 or more) above
//a base with n Nodes.
func LevelWidth(n Nodes, l Level) Nodes {
	if l < 0 {
		panic("can't see below")
	}
	for l > 0 {
		n = (n + 1) / 2
		l--
	}
	return n
}

//HashNumber gives each node in a tree an unique number from 0 up
func HashNumber(leafs Nodes, l Level, n Nodes) int64 {
	sum := Nodes(0)
	for i := Level(0); i < l; i++ {
		sum += LevelWidth(leafs, i)
	}
	return int64(sum + n)
}

//HashTreeSize is the total number of notes in a tree
func HashTreeSize(leafs Nodes) int64 {
	return HashNumber(leafs, Levels(leafs), 0)
}

//HashPosition uses hash HashNumber to tell you where you can
//put/get an inner hash in a byte array, without overlaps or unused space.
func HashPosition(leafs Nodes, l Level, n Nodes) int64 {
	return int64(HashNumber(leafs, l, n)) * int64(HashSize)
}

//SplitLength split the length of covered by an inner node in the tree to the
//length covered by it's childs. Such that:
//
//	b must be larger than LeafBlockSize (1024)
//	b = l + r
//	l is the largest possible power of 2, and l < b
//	r > 0
func SplitLength(b int64) (l, r int64) {
	if b <= LeafBlockSize {
		panic("can't split a leaf node")
	}
	mask := int64(highBitMask64(uint64(b)))
	if mask == b {
		l = b / 2
		r = l
	} else {
		l = mask
		r = b - l
	}
	return
}
