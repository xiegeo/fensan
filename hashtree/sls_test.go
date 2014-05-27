package hashtree

import (
	"testing"
)

type slsTest struct {
	from   Nodes
	to     Nodes
	width  Nodes
	expect [][2]Nodes
}

var slstestdata = []slsTest{
	{0, 0, 1, [][2]Nodes{{0, 0}}},
	{0, 0, 2, nil},
	{0, 1, 2, [][2]Nodes{{0, 1}}},
	{0, 2, 3, [][2]Nodes{{0, 2}}},
	{0, 2, 4, [][2]Nodes{{0, 1}}},
	{1, 2, 3, [][2]Nodes{{2, 2}}},
	{1, 2, 4, nil},
	{0, 3, 4, [][2]Nodes{{0, 3}}},
	{1, 6, 8, [][2]Nodes{{2, 3}, {4, 5}}},
	{1, 9, 10, [][2]Nodes{{2, 3}, {4, 7}, {8, 9}}},
	{2, 5, 10, [][2]Nodes{{2, 3}, {4, 5}}},
	{6, 9, 10, [][2]Nodes{{6, 7}, {8, 9}}},
	{10, 18, 20, [][2]Nodes{{10, 11}, {12, 15}, {16, 17}}},
	{1, 4, 5, [][2]Nodes{{2, 3}, {4, 4}}},
	{1, 4, 6, [][2]Nodes{{2, 3}}},
}

func TestSplitLocalSummable(t *testing.T) {
	for _, v := range slstestdata {
		sls := sls(v.from, v.to, v.width)
		if len(sls) != len(v.expect) {
			t.Errorf("%v got %v", v, sls)
		} else {
			for i := 0; i < len(sls); i++ {
				exp := v.expect[i]
				got := sls[i]
				if exp[0] != got[0] || exp[1] != got[1] {
					t.Errorf("part %v, got %v != exp %v, for test: %v", i, got, exp, v)
				}
			}
		}
	}
}
