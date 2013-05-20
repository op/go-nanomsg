// Go binding for nanomsg

package nanomsg

// #include <nanomsg/fanout.h>
import "C"

const (
	PUSH = Protocol(C.NN_PUSH)
	PULL = Protocol(C.NN_PULL)
)

type PushSocket struct {
	*Socket
}

// NewPushSocket creates a socket which is used to send messages to the
// cluster of load-balanced nodes. Receive operation is not implemented on
// this socket type.
func NewPushSocket() (*PushSocket, error) {
	socket, err := NewSocket(AF_SP, PUSH)
	return &PushSocket{socket}, err
}

type PullSocket struct {
	*Socket
}

// NewPullSocket creates a socket which is used to implement a node within
// a load-balanced cluster. It can be used to receive messages. Send
// operation is not implemented on this socket type. This socket can be
// connected to at most one peer.
func NewPullSocket() (*PullSocket, error) {
	socket, err := NewSocket(AF_SP, PULL)
	return &PullSocket{socket}, err
}
