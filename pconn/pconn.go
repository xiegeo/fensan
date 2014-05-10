/*
PConns are warpers on top of lower level network apis and complementing
algorithms (future: encoding, encryption, compression)

The primary purpose is to create a networking api for handling connections
between two end points, first based on TCP, that is easy to use to send data
(mostly in the Protocol Buffers format), that is also upgradable in the future
taking advantages of encryption and UDP without large api changes.

The preservation of message boundaries in PConns have several advantages vs streaming:
- Protocol Buffers does not self terminate, so that where one message end and the next start needs to be added before written to a stream such as TCP.
- Exposing message boundaries to encryption algorithms allows protection against truncation attacks. A Message is either processed in full, or not at all. To guarantee that a message is processed, an application level reply is still necessary.
- Easy of switching to a lossy transport such as UDP.

Features to manage multiple connections: a resource pool, will be somewhere else.
Features to create compatible connections with existing internet services will
probably not use this api.
*/
package pconn

import (
	"io"
)

//PConn is a generic network connection that preserves message boundaries, so
//that users don't need to worry about separating data in a stream.
//One Send function call will be matched by one Receive on the other end.
//
//Multiple goroutines may NOT invoke methods on a Conn simultaneously, unless supported.
type PConn interface {

	//Send sends a message. Send may block when buffer is full.
	//
	//Error Condtitions: those from net.Conn Write.
	//A good default error handling strategy is to close the connection.
	Send(msg []byte) error

	//Receive receives a message and have decoder process it.
	//Receive blocks untill decoder have finnshed processing the data unless
	//error is not nil.
	//
	//Error Conditions: those from net.Conn Read, or data corruption detected by
	//this PConn or the decoder.
	//A good default error handling strategy is to close the connection.
	Receive(decoder io.Writer) error

	//MaxMsgLength is the max length of a msg.
	//
	//Checked by the sender to crash on design errors, when violated a panic is
	//throw. Large message need to be breaked into smaller messages.
	//
	//Checked by the receiver to prevent large resource claims. Receive will error.
	//
	//Can be used by a PConn wapper (TODO) to know what size to splite large messages.
	MaxMsgLength() int

	//Close closes the connection. see net.Conn Close()
	Close() error
}

//ReceiveBytes get bytes from the Receive function of PConn
func ReceiveBytes(p PConn) ([]byte, error) {
	b := &bytesDecoder{}
	err := p.Receive(b)
	if err != nil {
		return nil, err
	}
	return b.buf, nil
}

type bytesDecoder struct {
	buf []byte
}

func (b *bytesDecoder) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}
