package store

import (
	"os"
	"testing"

	ht "github.com/xiegeo/fensan/hashtree"
)

type testPartHash struct {
	height ht.Level
	from   ht.Nodes
	length ht.Nodes
	count  ht.Nodes
}

func testSetUp(t *testing.T) (source MetaStore, part MetaStore) {
	sourceFolder := ".testSourceMetaStore"
	partFolder := ".testPartMetaStore"
	err := os.RemoveAll(sourceFolder)
	if err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll(partFolder)
	if err != nil {
		t.Fatal(err)
	}
	source, _ = OpenMetaStore(sourceFolder)
	part, _ = OpenMetaStore(partFolder)
	return
}
