// Go binding for nanomsg

package nanomsg

// #include <nanomsg/tcp.h>
import "C"

// TcpNoDelay returns the current value of TCP no delay.
func (s *Socket) TcpNoDelay() (bool, error) {
	noDelay, err := s.SockOptInt(C.NN_TCP, C.NN_TCP_NODELAY)
	return noDelay != 0, err
}

// SetTcpNoDelay controls whether the operating system should delay packet
// transmission in hopes of sending fewer packets (Nagle's algorithm).
func (s *Socket) SetTcpNoDelay(noDelay bool) error {
	var value int
	if noDelay {
		value = 1
	}
	return s.SetSockOptInt(C.NN_TCP, C.NN_TCP_NODELAY, value)
}
