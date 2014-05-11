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
//Use SendBytes and ReceiveBytes for a []byte interface
//
//Multiple goroutines may NOT invoke methods on a Conn simultaneously, unless supported.
type PConn interface {

	//Sender sends a message written to the returned writer. The message must end
	//by Close, else nothing should happen on the other end.
	//Sender's Write and Close may block when network buffer is full.
	//
	//Unless supported, sender can only be called again after the last sender is closed.
	//Sender may rereturn a old writer, previously closed, to send a new message.
	//
	//Error Condtitions: those from net.Conn Write.
	//A good default error handling strategy is to close the connection (not the sender).
	Sender() io.WriteCloser

	//Receive receives a message as a reader. read EOF is the end of message.
	//Receive may block, or return reader like a future.
	//
	//Unless supported, Receive should not be called again untill read report EOF. ReceiveInWriter
	//and ReceiveBytes does this for you by blocking.
	//
	//All Errors are reported by read. Any error, other then EOF means all actions
	//taken by reader must be reverted for security purposes.
	//Error Conditions: those from net.Conn Read, or data corruption detected by
	//this PConn.
	//A good default error handling strategy is to close the connection.
	Receiver() io.Reader

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

//SendBytes sends a message using Sender
func SendBytes(p PConn, msg []byte) error {
	s := p.Sender()
	_, err := s.Write(msg)
	if err != nil {
		return err
	}
	return s.Close()
}

//ReceiveInWriter receives a full message and have decoder process it.
//Receive blocks untill decoder have finnshed processing the data unless
//error is not nil.
func ReceiveInWriter(p PConn, decoder io.Writer) error {
	r := p.Receiver()
	_, err := io.Copy(decoder, r)
	return err
}

//ReceiveBytes get bytes from the Receive function of PConn
//Returns all bytes of a message will nil error, or nil message with some error
func ReceiveBytes(p PConn) ([]byte, error) {
	b := &bytesDecoder{}
	err := ReceiveInWriter(p, b)
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
