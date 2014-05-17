package hashtree

import (
	"testing"
)

type treeLevels struct {
	n      Nodes
	levels Level
}

var expectedTreeLevels = []treeLevels{
	{1, 1},
	{2, 2},
	{3, 3}, {4, 3},
	{5, 4}, {8, 4},
	{9, 5}, {16, 5},
	{17, 6},
}

func TestTreeLevels(t *testing.T) {
	for i := 0; i < len(expectedTreeLevels); i++ {
		e := expectedTreeLevels[i]
		if Levels(e.n) != e.levels {
			t.Fatalf("Levels(%d) = %d want %d", e.n, Levels(e.n), e.levels)
		}
	}
}

type levelWidth struct {
	n     Nodes
	level Level
	width Nodes
}

var expectedLevelWidth = []levelWidth{
	//{1, -1, 0}, {2, -1, 0}, {3, -1, 0}, //special case for level < 0, no nodes exist here so the width is 0

	{1, 0, 1},

	{2, 0, 2},
	{2, 1, 1},

	{3, 0, 3}, {4, 0, 4},
	{3, 1, 2}, {4, 1, 2},
	{3, 2, 1}, {4, 2, 1},

	{5, 0, 5}, {6, 0, 6}, {7, 0, 7}, {8, 0, 8},
	{5, 1, 3}, {6, 1, 3}, {7, 1, 4}, {8, 1, 4},
	{5, 2, 2}, {6, 2, 2}, {7, 2, 2}, {8, 2, 2},
	{5, 3, 1}, {6, 3, 1}, {7, 3, 1}, {8, 3, 1},
}

func TestLevelWidth(t *testing.T) {
	for i := 0; i < len(expectedLevelWidth); i++ {
		e := expectedLevelWidth[i]
		if LevelWidth(e.n, e.level) != e.width {
			t.Fatalf("LevelWidth(%d,%d) = %d want %d", e.n, e.level, LevelWidth(e.n, e.level), e.width)
		}
	}
}

type nodeChildLength struct{ b, l, r int64 }

var expectedCL = []nodeChildLength{
	{1025, 1024, 1}, {2000, 1024, 2000 - 1024}, {3000, 2048, 3000 - 2048},
	{2048, 1024, 1024},
}

func TestSplitLength(t *testing.T) {
	for _, ex := range expectedCL {
		l, r := SplitLength(ex.b)
		if l != ex.l || r != ex.r {
			t.Fatalf("test failed")
		}
	}
}

func TestInnerHashListener(t *testing.T) {
	testInnerHashListener([][]int32{
		{0},
	}, t)
	testInnerHashListener([][]int32{
		{0, 1},
		{-1},
	}, t)
	testInnerHashListener([][]int32{
		{0, 1, 1},
		{-1, 1},
		{-2},
	}, t)
	testInnerHashListener([][]int32{
		{0, 1, 1, 3},
		{-1, -2},
		{1},
	}, t)
	testInnerHashListener([][]int32{
		{1, 1, 2, 3, 5, 8},
		{0, -1, -3},
		{1, -3},
		{4},
	}, t)
}
func testInnerHashListener(inner [][]int32, t *testing.T) {
	//t.Log(inner)
	listener := func(l Level, i Nodes, hash, left, right *H256) {
		//t.Log(l, i, hash, left, right)
		/*defer func() {
			if r := recover(); r != nil {
				t.Fatalf("error:%s, at Level:%d, Node:%d ", r, l, i)
			}
		}()*/
		h := int32(hash[0])
		if inner[l][i] != h {
			if inner[l][i] == h+2000 {
				t.Fatalf("Level:%d, Node:%d was repeated", l, i)
			}
			t.Fatalf("Level:%d, Node:%d, hash:%d, should be %d", l, i, h, inner[l][i])
		}
		inner[l][i] += 2000 //mark heard

		//check child
		if l == 0 {
			if left != nil || right != nil {
				t.Fatal("childs must be nil on base level")
			}
		} else {
			if right == nil {
				lc := int32(left[0])
				if h != lc {
					t.Fatal("promoted must idential")
				}
				if lc+2000 != inner[l-1][i*2] {
					t.Fatal("child miss-match")
				}
			} else {
				lc := int32(left[0]) + 2000
				rc := int32(right[0]) + 2000
				if lc != inner[l-1][i*2] {
					t.Fatalf("child miss-match:%v!=%v@%v,%v", lc, inner[l-1][i*2], l, i)
				}
				if rc != inner[l-1][i*2+1] {
					t.Fatal("child miss-match")
				}
			}
		}
	}
	c := NewTree2(NoPad32bytes, minus).(*treeDigest)
	c.SetInnerHashListener(listener)
	for _, n := range inner[0] {
		data := H256{uint32(n)}
		c.Write(data.ToBytes())
	}
	c.Sum(nil)

	for l, array := range inner {
		for i, n := range array {
			if n < 1000 {
				t.Fatalf("Level:%d, Node:%d was not heard", l, i)
			}
		}
	}
}
