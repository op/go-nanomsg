// Go binding for nanomsg

package nanomsg

// #include <nanomsg/pubsub.h>
import "C"

const (
	PUB = Protocol(C.NN_PUB)
	SUB = Protocol(C.NN_SUB)
)

type SubSocket struct {
	*Socket
}

// NewSubSocket creates a new socket which receives messages from the
// publisher. Only messages that the socket is subscribed to are received. When
// the socket is created there are no subscriptions and thus no messages will
// be received. Send operation is not defined on this socket. The socket can
// be connected to at most one peer.
func NewSubSocket() (*SubSocket, error) {
	socket, err := NewSocket(AF_SP, SUB)
	return &SubSocket{socket}, err
}

// Subscribe subscribes to a particular topic.
func (sub *SubSocket) Subscribe(topic string) error {
	return sub.SetSockOptString(C.NN_SUB, C.NN_SUB_SUBSCRIBE, topic)
}

// Unsubscribe unsubscribes from a particular topic.
func (sub *SubSocket) Unsubscribe(topic string) error {
	return sub.SetSockOptString(C.NN_SUB, C.NN_SUB_UNSUBSCRIBE, topic)
}

type PubSocket struct {
	*Socket
}

// NewPubSocket creates a new socket which is used to distribute messages to
// multiple destinations. Receive operation is not defined.
func NewPubSocket() (*PubSocket, error) {
	socket, err := NewSocket(AF_SP, PUB)
	return &PubSocket{socket}, err
}
