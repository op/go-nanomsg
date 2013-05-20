// Go binding for nanomsg

package nanomsg

// #include <nanomsg/pair.h>
import "C"

const (
	PAIR = Protocol(C.NN_PAIR)
)

type PairSocket struct {
	*Socket
}

// NewPairSocket creates a socket for communication with exactly one peer. Each
// party can send messages at any time. If the peer is not available or send
// buffer is full, subsequent calls to Send will block until it’s possible
// to send the message.
func NewPairSocket() (*PairSocket, error) {
	socket, err := NewSocket(AF_SP, PAIR)
	return &PairSocket{socket}, err
}
