package hashtree

import (
	"fmt"
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

//HashTreeSize is the total number of nodes in a tree
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

//Split inner hashes based on highest derivable ancestors
func SplitLocalSummable(hashes []byte, hashSize int, levelWidth Nodes, off Nodes) [][]byte {
	length := Nodes(len(hashes) / hashSize)
	from := off
	to := from + length - 1
	ranges := slsUntrusted(from, to, levelWidth)
	if ranges == nil {
		return nil
	}
	return split(hashes, hashSize, off, ranges)
}

func split(b []byte, hashSize int, off Nodes, ranges [][2]Nodes) [][]byte {
	r := make([][]byte, len(ranges))
	for i, v := range ranges {
		fb := int((v[0] + off)) * hashSize
		tb := int((v[1] + off)) * hashSize
		r[i] = b[fb:tb]
	}
	return r
}

// the input maybe gernerated by an adversary, return nil instead of panic that check coding errors
func slsUntrusted(from, to, width Nodes) [][2]Nodes {
	if from > to || to >= width {
		return nil
	} else {
		return sls(from, to, width)
	}
}

func sls(from, to, width Nodes) [][2]Nodes {
	from = (from + 1) / 2 * 2
	if from > to || to >= width {
		panic(fmt.Sprintf("from:%v, to:%v, width:%v", from, to, width))
	}

	if from == to {
		//there souldn't be any singles, unless it is the last one and even
		if from == width-1 && from%2 == 0 {
			return [][2]Nodes{{from, to}}
		}
		return nil
	}
	if from == 0 {
		dev := Nodes(highBitMask(uint32(to + 1)))
		//log.Println(from,to,width,dev);
		if to == width-1 || to == dev-1 {
			return [][2]Nodes{{from, to}}
		}
		return mergeR(sls(from, dev-1, dev), shiftsls(sls(0, to-dev, width-dev), dev))
	} else {
		dev := Nodes(highBitMask(uint32(width - 1)))
		//log.Println(from,to,width,dev);
		if from < dev {
			if to < dev {
				return sls(from, to, dev)
			} else {
				return mergeR(sls(from, dev-1, dev), shiftsls(sls(0, to-dev, width-dev), dev))
			}
		} else {
			return shiftsls(sls(from-dev, to-dev, width-dev), dev)
		}
	}
}

func mergeR(a [][2]Nodes, b [][2]Nodes) [][2]Nodes {
	result := make([][2]Nodes, len(a)+len(b))
	copy(result, a)
	copy(result[len(a):], b)
	return result
}

func shiftsls(sls [][2]Nodes, delta Nodes) [][2]Nodes {
	for i := 0; i < len(sls); i++ {
		sls[i][0] += delta
		sls[i][1] += delta
	}
	return sls
}
