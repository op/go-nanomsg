// Go binding for nanomsg

package nanomsg

// #include <nanomsg/reqrep.h>
import "C"

import (
	"time"
)

const (
	REQ = Protocol(C.NN_REQ)
	REP = Protocol(C.NN_REP)
)

type ReqSocket struct {
	*Socket
}

// NewReqSocket creates a request socket used to implement the client
// application that sends requests and receives replies.
func NewReqSocket() (*ReqSocket, error) {
	socket, err := NewSocket(AF_SP, REQ)
	return &ReqSocket{socket}, err
}

// ResendInterval returns the resend interval. If reply is not received in
// specified amount of time, the request will be automatically resent. Default
// value is 1 minute.
func (req *ReqSocket) ResendInterval() (time.Duration, error) {
	return req.Socket.SockOptDuration(C.NN_REQ, C.NN_REQ_RESEND_IVL, time.Millisecond)
}

// SetResendInterval sets the resend interval for requests.
func (req *ReqSocket) SetResendInterval(interval time.Duration) error {
	return req.Socket.SetSockOptDuration(C.NN_REQ, C.NN_REQ_RESEND_IVL, time.Millisecond, interval)
}

type RepSocket struct {
	*Socket
}

// NewRepSocket creates a reply socket used to implement the stateless worker
// that receives requests and sends replies.
func NewRepSocket() (*RepSocket, error) {
	socket, err := NewSocket(AF_SP, REP)
	return &RepSocket{socket}, err
}
