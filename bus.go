// Go binding for nanomsg

package nanomsg

// #include <nanomsg/bus.h>
import "C"

const (
	BUS = Protocol(C.NN_BUS)
)

type BusSocket struct {
	*Socket
}

// NewBusSocket creates a socket where sent messages are distributed to all
// nodes in the topology. Incoming messages from all other nodes in the
// topology are fair-queued in the socket.
func NewBusSocket() (*BusSocket, error) {
	socket, err := NewSocket(AF_SP, BUS)
	return &BusSocket{socket}, err
}
