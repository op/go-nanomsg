// Go binding for nanomsg

package nanomsg

// #include <nanomsg/tcp.h>
import "C"

// TCPNoDelay returns the current value of TCP no delay.
func (s *Socket) TCPNoDelay() (bool, error) {
	return s.SockOptBool(C.NN_TCP, C.NN_TCP_NODELAY)
}

// SetTCPNoDelay controls whether the operating system should delay packet
// transmission in hopes of sending fewer packets (Nagle's algorithm).
func (s *Socket) SetTCPNoDelay(noDelay bool) error {
	return s.SetSockOptBool(C.NN_TCP, C.NN_TCP_NODELAY, noDelay)
}
