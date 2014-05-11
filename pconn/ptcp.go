package pconn

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

//PTCP implements PConn over net.TCPConn
type PTCP struct {
	tcp                 *net.TCPConn
	writeBuf            *pTCPSender
	readBuf             *bufio.Reader
	lastReader          *io.LimitedReader //to check that the last read is finished
	debugSendCounter    byte              //not really needed just an extra check on the data
	debugReceiveCounter byte
}

func NewPTCP(c *net.TCPConn) *PTCP {
	return &PTCP{c, nil, bufio.NewReader(c), nil, 0, 0}
}

func (p *PTCP) Sender() io.WriteCloser {
	return p.renewWriteBuf()
}

func (p *PTCP) Receiver() io.Reader {
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

	h, err := p.readBuf.ReadByte()
	if err != nil {
		return er(err)
	}
	l, err := p.readBuf.ReadByte()
	if err != nil {
		return er(err)
	}
	length := int(h)<<8 + int(l)
	if length > p.MaxMsgLength() {
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
		panic("errorReader does not support other errors that also reads")
	}
	return errorReader{err}
}

func erf(format string, a ...interface{}) errorReader {
	return er(fmt.Errorf(format, a...))
}

func (e errorReader) Read(ignore []byte) (int, error) {
	return 0, e
}
