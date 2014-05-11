package pconn

//pTCPSender is a reusable WriteCloser returned by PTCP.Sender
type pTCPSender struct {
	conn           *PTCP
	inUse          bool
	buf            []byte
	payLoadCounter int
}

const pTCPSenderLengthSize = 2
const pTCPSenderLengthOffset = 1

func (p *PTCP) renewWriteBuf() *pTCPSender {
	s := p.writeBuf
	if s == nil {
		s = &pTCPSender{
			conn:           p,
			inUse:          true,
			buf:            make([]byte, 0, p.MaxMsgLength()+pTCPSenderLengthOffset+pTCPSenderLengthSize),
			payLoadCounter: 0,
		}
		p.writeBuf = s
	} else {
		if s.inUse {
			panic("Must close last Sender before calling Sender again.")
		}
		s.buf = s.buf[:0]
		s.inUse = true
		s.payLoadCounter = 0
	}
	s.buf = append(s.buf, s.conn.debugSendCounter)
	s.conn.debugSendCounter++
	s.buf = append(s.buf, make([]byte, 2)...)
	return s
}

func (s *pTCPSender) Write(p []byte) (n int, err error) {
	s.payLoadCounter += len(p)
	if s.payLoadCounter > s.conn.MaxMsgLength() {
		panic("send msg is too long")
	}
	s.buf = append(s.buf, p...)
	return len(p), nil
}

func (s *pTCPSender) Close() error {
	if !s.inUse {
		panic("already closed")
	}
	s.inUse = false
	n := int16(s.payLoadCounter)
	s.buf[1], s.buf[2] = byte(n>>8), byte(n&0xff)
	//fmt.Print(s.buf)
	_, err := s.conn.tcp.Write(s.buf) //one write call, avoid sending many packets when on no delay.
	return err
}
