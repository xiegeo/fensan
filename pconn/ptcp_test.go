package pconn

import (
	"bytes"
	"net"
	"testing"
)

func TestPTCP(t *testing.T) {
	protocal := "tcp"
	serverAddr, err := net.ResolveTCPAddr(protocal, "localhost:5656")
	assertNil(err)
	listener, err := net.ListenTCP(protocal, serverAddr)
	assertNil(err)

	data := [][]byte{{5}, {6, 7, 8}, make([]byte, 4096)}

	go func() {
		request, err := net.DialTCP(protocal, nil, serverAddr)
		assertNil(err)
		rc := NewPTCP(request)
		for _, d := range data {
			err = SendBytes(rc, d)
			assertNil(err)
		}
		//try send bad data
		request.Write([]byte{99})

		rc.Close()
	}()

	accept, err := listener.AcceptTCP()
	assertNil(err)
	ac := NewPTCP(accept)
	for _, d := range data {
		got, err := ReceiveBytes(ac)
		assertNil(err)
		if !bytes.Equal(got, d) {
			t.Error("send:", d, " but received:", got)
		}
	}
	//try get bad data
	got, err := ReceiveBytes(ac)
	if len(got) != 0 || err == nil {
		t.Error("expecting bad data, but got:", got, " with error:", err)
	}
}

func assertNil(e error) {
	if e != nil {
		panic(e)
	}
}
