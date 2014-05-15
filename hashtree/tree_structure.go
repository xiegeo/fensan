package hashtree

import (
	"math"
)

//Level is the depth of tree counted from bottom up.
//The lowest level is 0.
//Level is signed to allow representation of deltas.
type Level int

//Bytes is the size in bytes.
//Bytes is signed to allow representation of deltas.
//Bytes is 64 bits like in the file io apis
type Bytes int64

//Nodes is the number of Nodes in a level,
//or the index of a node in a level.
//Nodes is signed to allow representation of deltas.
type Nodes int

//Levels return the number of Levels (1 or more) of a hash
//tree with n leaf nodes
func Levels(n Nodes) Level {
	return Level(math.Ilogb(float64(n*2-1)) + 1)
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
func HashPosition(leafs Nodes, l Level, n Nodes) Bytes {
	return Bytes(HashNumber(leafs, l, n)) * Bytes(HashSize)
}
