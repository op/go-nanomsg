// Go binding for nanomsg

package nanomsg

// #include <nanomsg/nn.h>
import "C"

import (
	"time"
)

// Poller is used to poll a set of sockets for readability and/or writability.
type Poller struct {
	fds []C.struct_nn_pollfd
}

// Add puts the given socket into the poller to check when it's available for
// sending or receiving. Use the returned PollItem to check what state the
// socket is in or to modify what events to wait for.
func (p *Poller) Add(s *Socket, recv, send bool) *PollItem {
	var fd C.struct_nn_pollfd
	fd.fd = s.socket
	pi := &PollItem{p, len(p.fds)}
	p.fds = append(p.fds, fd)
	pi.PollRecv(recv)
	pi.PollSend(send)
	return pi
}

// Poll returns as soon as any of the sockets are available for sending and/or
// receiving, depending on how the poll item is setup. The timeout is used to
// specify how long the function should block if there are no events.
//
// This function returns the number of events and error. If the poller timed
// out before any event was received, the number of events will be 0.
func (p *Poller) Poll(timeout time.Duration) (int, error) {
	t := C.int(timeout / time.Millisecond)
	rc, err := C.nn_poll(&p.fds[0], C.int(len(p.fds)), t)
	if rc == -1 {
		return 0, nnError(err)
	}
	return int(rc), nil
}

// PollItem represents a socket and what events to poll for.
type PollItem struct {
	poller *Poller
	index  int
}

// PollRecv is used to specify if the poller should return as soon as the
// socket is ready to receive data from.
func (pi *PollItem) PollRecv(recv bool) {
	if recv {
		pi.poller.fds[pi.index].events |= C.NN_POLLIN
	} else {
		pi.poller.fds[pi.index].events ^= C.NN_POLLIN
	}
}

// PollSend is used to specify if the poller should return as soon as the
// socket is ready to send data on.
func (pi *PollItem) PollSend(send bool) {
	if send {
		pi.poller.fds[pi.index].events |= C.NN_POLLOUT
	} else {
		pi.poller.fds[pi.index].events ^= C.NN_POLLOUT
	}
}

// CanRecv returns true if the socket is ready to receive data from.
func (pi *PollItem) CanRecv() bool {
	return (pi.poller.fds[pi.index].revents & C.NN_POLLIN) == C.NN_POLLIN
}

// CanSend returns true if the socket is ready to send data on.
func (pi *PollItem) CanSend() bool {
	return (pi.poller.fds[pi.index].revents & C.NN_POLLOUT) == C.NN_POLLOUT
}
