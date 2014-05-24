package bitset

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

type randomFlushingBitSet struct {
	FileBackedBitSet
	r *rand.Rand
}

func (b *randomFlushingBitSet) Set(i int) {
	b.setAs(i, true)
	if b.r.Float32() < 0.02 {
		b.Flush()
	}
}

func (b *randomFlushingBitSet) Unset(i int) {
	b.setAs(i, false)
	if b.r.Float32() < 0.02 {
		b.Flush()
	}
}

func testFileBacked(cap int, t *testing.T) {
	fileName := fmt.Sprintf(".testfile_%v", cap)
	nfb := OpenFileBacked(fileName, cap)
	s := &randomFlushingBitSet{*nfb, rand.New(rand.NewSource(0))}
	checkAll(t, s, cap)
	if s.Capacity() != cap {
		t.Fatalf("capacity should be %v but returns %v", cap, s.Capacity())
	}
	if CHECK_INTEX {
		tryOutSide(s, -1, t)
		tryOutSide(s, cap, t)
	}
	s.Close()
	err := os.Remove(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileSet0(t *testing.T)   { testFileBacked(0, t) }
func TestFileSet1(t *testing.T)   { testFileBacked(1, t) }
func TestFileSet2(t *testing.T)   { testFileBacked(2, t) }
func TestFileSet7(t *testing.T)   { testFileBacked(7, t) }
func TestFileSet8(t *testing.T)   { testFileBacked(8, t) }
func TestFileSet9(t *testing.T)   { testFileBacked(9, t) }
func TestFileSet127(t *testing.T) { testFileBacked(127, t) }
func TestFileSet128(t *testing.T) { testFileBacked(128, t) }
func TestFileSet129(t *testing.T) { testFileBacked(129, t) }
