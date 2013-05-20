// Go binding for nanomsg

package nanomsg

// #include <nanomsg/fanin.h>
import "C"

const (
	SOURCE = Protocol(C.NN_SOURCE)
	SINK   = Protocol(C.NN_SINK)
)

type SourceSocket struct {
	*Socket
}

// NewSourceSocket creates a socket which allows to send messages to the
// central sink. Receive operation is not implemented on this socket type.
// This socket can be connected to at most one peer.
func NewSourceSocket() (*SourceSocket, error) {
	socket, err := NewSocket(AF_SP, SOURCE)
	return &SourceSocket{socket}, err
}

type SinkSocket struct {
	*Socket
}

// NewSinkSocket creates a socket which receives the messages from multiple
// sources. Send operation is not defined on this socket type.
func NewSinkSocket() (*SinkSocket, error) {
	socket, err := NewSocket(AF_SP, SINK)
	return &SinkSocket{socket}, err
}
