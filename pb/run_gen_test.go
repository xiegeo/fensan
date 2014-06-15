package pb

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

//test if the code generater had run after last edit of .proto file
func TestRunGen(t *testing.T) {
	ls, _ := ioutil.ReadDir(".")
	protos := make(map[string]os.FileInfo)
	pbgos := make(map[string]os.FileInfo)
	for _, f := range ls {
		n := f.Name()
		if strings.HasSuffix(n, ".proto") {
			protos[strings.TrimSuffix(n, ".proto")] = f
		} else if strings.HasSuffix(n, ".pb.go") {
			pbgos[strings.TrimSuffix(n, ".pb.go")] = f
		}
	}
	if len(protos) != len(pbgos) {
		t.Error("Please do a clean regeneration of proto code")
	}
	for n, info := range protos {
		if info.ModTime().After(pbgos[n].ModTime()) {
			t.Error("Please do a clean regeneration of proto code")
		}
	}
}
