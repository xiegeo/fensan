package main

import (
	"fmt"
	"os/exec"
)

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
	err := doCmd(exec.Command("go", "test", "-v"))
	if err != nil {
		doCmd(exec.Command("protoc", "--gogo_out=.", "*.proto"))
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
		fmt.Println(out)
		fmt.Printf("%s\n", out)
		fmt.Printf("error:%v\n", err)
	}
	return err
}
