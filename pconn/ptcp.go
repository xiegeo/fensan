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
	debugSendCounter    byte //not really needed just an extra check on the data
	debugReceiveCounter byte
	buf                 []byte
}

func NewPTCP(c *net.TCPConn) *PTCP {
	return &PTCP{c, bufio.NewReader(c), 0, 0, nil}
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

func (p *PTCP) Receive(decoder io.Writer) error {
	count, err := p.readBuf.ReadByte()
	if err != nil {
		return err
	}
	if count != p.debugReceiveCounter {
		return fmt.Errorf("data corrupted: debug counter mismatch:%v,%v", count, p.debugReceiveCounter)
	}
	p.debugReceiveCounter++

	length, err := binary.ReadUvarint(p.readBuf)
	if err != nil {
		return err
	}
	if length > uint64(p.MaxMsgLength()) {
		return fmt.Errorf("data corrupted: data is too long")
	}

	pReader := io.LimitReader(p.readBuf, int64(length))

	written, err := io.Copy(decoder, pReader)
	if err != nil {
		return err
	}
	if written != int64(length) {
		panic("must never happen without errs")
	}
	return nil
}

func (p *PTCP) MaxMsgLength() int {
	return 4096
}

func (p *PTCP) Close() error {
	return p.tcp.Close()
}
