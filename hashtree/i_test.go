package hashtree

import "testing"

type iTestdata struct {
	b Bytes
	n Nodes
}

var iExpected = []iTestdata{
	{0, 1}, {1, 1}, {1024, 1},
	{1025, 2}, {2048, 2},
	{2049, 3},
}

func TestI(t *testing.T) {
	for _, d := range iExpected {
		got := I.Nodes(d.b)
		if got != d.n {
			t.Fatal(d.b, " bytes should gave ", d.n, "nodes but got ", got)
		}
	}
}
