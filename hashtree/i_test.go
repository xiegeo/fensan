package hashtree

import "testing"

type iTestdata struct {
	b Bytes
	n Nodes
}

func TestI(t *testing.T) {
	iExpected := []iTestdata{
		{0, 1}, {1, 1}, {1024, 1},
		{1025, 2}, {2048, 2},
		{2049, 3},
	}
	for _, d := range iExpected {
		got := I.Nodes(d.b)
		if got != d.n {
			t.Fatal(d.b, " bytes should gave ", d.n, "nodes but got ", got)
		}
	}
}

type iTestHashNumberingData struct {
	leafs Nodes
	total int64
}

func TestHashNumbering(t *testing.T) {
	testData := []iTestHashNumberingData{
		{1, 1}, {2, 3}, {3, 6}, {4, 7}, {5, 11},
		{6, 12}, {7, 14}, {8, 15}, {9, 20}, {10, 21},
	}
	for _, d := range testData {
		gotTotal := HashTreeSize(d.leafs)
		if gotTotal != d.total {
			t.Fatal(d.leafs, " leafs should produce ", d.total, " total nodes but got ", gotTotal)
		}
		tSize := d.total * HashSize
		toFill := make([]byte, tSize)

		levels := Levels(d.leafs)
		for i := Level(0); i < levels; i++ {
			width := LevelWidth(d.leafs, i)
			for j := Nodes(0); j < width; j++ {
				p := HashPosition(d.leafs, i, j)
				if toFill[p] != 0 {
					t.Fatal("position already taken")
				}
				toFill[p] = 1
			}
		}
		for i := int64(0); i < tSize; i++ {
			if i%HashSize == 0 {
				if toFill[i] == 0 {
					t.Fatal("position not taken")
				}
			} else {
				if toFill[i] == 1 {
					t.Fatal("position taken shifted out of place")
				}
			}
		}
	}
}
