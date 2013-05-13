// Go binding for nanomsg

package reqrep

// #include <nanomsg/reqrep.h>
import "C"

import (
	"time"

	//"github.com/op/go-nanomsg"
	".."
)

type ReqSocket struct {
	*nanomsg.Socket
}

// NewReqSocket creates a request socket used to implement the client
// application that sends requests and receives replies.
func NewReqSocket() (*ReqSocket, error) {
	socket, err := nanomsg.NewSocket(nanomsg.SP, C.NN_REQ)
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
	ivlMs := int(linger / time.Millisecond)
	return s.SetSockOptInt(C.NN_REQ, C.NN_REQ_RESEND_IVL, lingerMs)
}

type RepSocket struct {
	*nanomsg.Socket
}

// NewRepSocket creates a reply socket used to implement the stateless worker
// that receives requests and sends replies.
func NewRepSocket() (*RepSocket, error) {
	socket, err := nanomsg.NewSocket(nanomsg.SP, C.NN_REP)
	return &RepSocket{socket}, err
}
