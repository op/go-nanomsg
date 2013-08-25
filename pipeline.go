// Go binding for nanomsg

package nanomsg

// #include <nanomsg/pipeline.h>
import "C"

const (
	PUSH = Protocol(C.NN_PUSH)
	PULL = Protocol(C.NN_PULL)
)

type PushSocket struct {
	*Socket
}

// NewPushSocket creates a socket which is used to send messages to a cluster
// of load-balanced nodes. Receive operation is not implemented on this socket
// type.
func NewPushSocket() (*PushSocket, error) {
	socket, err := NewSocket(AF_SP, PUSH)
	return &PushSocket{socket}, err
}

type PullSocket struct {
	*Socket
}

// NewPullSocket creates a socket which is used to receive a message from a
// cluster of nodes. Send operation is not implemented on this socket type.
func NewPullSocket() (*PullSocket, error) {
	socket, err := NewSocket(AF_SP, PULL)
	return &PullSocket{socket}, err
}
