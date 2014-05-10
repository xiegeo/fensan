package pconn

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type PTCP struct {
	tcp                 *net.TCPConn
	readBuf             *bufio.Reader
	lastReader          *io.LimitedReader //to check that the last read is finished
	debugSendCounter    byte              //not really needed just an extra check on the data
	debugReceiveCounter byte
	buf                 []byte
}

func NewPTCP(c *net.TCPConn) *PTCP {
	return &PTCP{c, bufio.NewReader(c), nil, 0, 0, nil}
}

const maxSendLengthIncrease = 11 //10 for varint + 1 for debug_counter

func (p *PTCP) Send(msg []byte) error {
	length := len(msg)
	if length > p.MaxMsgLength() {
		panic("send msg is too long")
	}
	buf := p.buf
	if buf == nil {
		buf = make([]byte, p.MaxMsgLength()+maxSendLengthIncrease)
	}
	buf = buf[:cap(buf)]
	size := binary.PutUvarint(buf, uint64(length))
	copy(buf[1:size+1], buf[:size])
	buf[0] = p.debugSendCounter
	p.debugSendCounter++
	copy(buf[size+1:], msg)
	buf = buf[:1+size+length]
	//fmt.Print(buf)
	_, err := p.tcp.Write(buf) //one write call, avoid sending many packets when on no delay.
	p.buf = buf
	return err
}

func (p *PTCP) Receive() io.Reader {
	if p.lastReader != nil && p.lastReader.N > 0 {
		panic("Receive is not ready for reused")
	}
	count, err := p.readBuf.ReadByte()
	if err != nil {
		return er(err)
	}
	if count != p.debugReceiveCounter {
		return erf("data corrupted: debug counter mismatch:%v,%v", count, p.debugReceiveCounter)
	}
	p.debugReceiveCounter++

	length, err := binary.ReadUvarint(p.readBuf)
	if err != nil {
		return er(err)
	}
	if length > uint64(p.MaxMsgLength()) {
		return erf("data corrupted: data is too long")
	}

	p.lastReader = io.LimitReader(p.readBuf, int64(length)).(*io.LimitedReader)
	return p.lastReader
}

func (p *PTCP) MaxMsgLength() int {
	return 4096
}

func (p *PTCP) Close() error {
	return p.tcp.Close()
}

type errorReader struct {
	error
}

//turns an error into a reader that just return this error on read
func er(err error) errorReader {
	reader, ok := err.(errorReader)
	if ok {
		return reader
	}
	_, notOk := err.(io.Reader)
	if notOk {
		panic("errorReader does not support other error that also reads")
	}
	return errorReader{err}
}

func erf(format string, a ...interface{}) errorReader {
	return er(fmt.Errorf(format, a...))
}

func (e errorReader) Read(ignore []byte) (int, error) {
	return 0, e
}
