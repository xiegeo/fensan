package main

import (
	"fmt"
	"os/exec"

	"github.com/xiegeo/fensan/bitset"
	"github.com/xiegeo/fensan/hashtree"
	"github.com/xiegeo/fensan/pb"
	"github.com/xiegeo/fensan/pconn"
	"github.com/xiegeo/fensan/store"
)

//make sure go get gets every sub package
var _ = bitset.CHECK_INTEX
var _ = hashtree.HashSize
var _ = &pb.StaticId{}
var _ = pconn.SendBytes
var _ = store.FileNone

func main() {
	buildProtoBuf()
	testCode("bitset")
	testCode("hashtree")
	testCode("pb")
	testCode("pconn")
	testCode("store")
	fmt.Println("\n\ndone all builds and tests")
}

func buildProtoBuf() {
	dir = "pb"
	defer func() { dir = "" }()
	err := doHiddenCmd(exec.Command("go", "test", "-v"))
	if err != nil {
		fmt.Println("rebuilding .pb.go files")
		err := doCmd(exec.Command("protoc", "--gogo_out=.", "*.proto"))
		if err != nil {
			panic("can't rebuild, see code.google.com/p/gogoprotobuf/")
		}
		if err == nil {
			fmt.Println("rebuilded .pb.go files ")
		}
	}
}

func testCode(packageName string) {
	dir = packageName
	defer func() { dir = "" }()
	noErr(doCmd(exec.Command("go", "test", "-v")))
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

var dir = ""

func doCmd(cmd *exec.Cmd) error {
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd.Path)
		fmt.Println(cmd.Args)
		fmt.Printf("%s\n", out)
		fmt.Printf("error:%v\n", err)
	}
	return err
}

func doHiddenCmd(cmd *exec.Cmd) error {
	cmd.Dir = dir
	_, err := cmd.CombinedOutput()
	return err
}
