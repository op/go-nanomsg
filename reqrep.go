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

// ResendIvl returns the resend interval. If reply is not received in specified
// amount of time, the request will be automatically resent. Default value is 1
// minute.
func (req *ReqSocket) ResendIvl() (time.Duration, error) {
	ivl, err := req.Socket.GetSockOptInt(C.NN_REQ, C.NN_REQ_RESEND_IVL)
	return time.Duration(ivl) * time.Millisecond, err
}

// SetResendIvl sets the resend interval for requests.
func (req *ReqSocket) SetResendIvl(ivl time.Duration) error {
	ivlMs := int(ivl / time.Millisecond)
	return req.Socket.SetSockOptInt(C.NN_REQ, C.NN_REQ_RESEND_IVL, ivlMs)
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
