// Go binding for nanomsg

package nanomsg

// #include <nanomsg/survey.h>
import "C"

import (
	"time"
)

const (
	SURVEYOR   = Protocol(C.NN_SURVEYOR)
	RESPONDENT = Protocol(C.NN_RESPONDENT)
)

type SurveyorSocket struct {
	*Socket
}

// NewSurveyorSocket creates a socket used to send the survey.
// The survey is delivered to all the connected respondents. Once
// the query is sent, the socket can be used to receive the
// responses. When the survey deadline expires, receive will
// return ETIMEDOUT error.
func NewSurveyorSocket() (*SurveyorSocket, error) {
	socket, err := NewSocket(AF_SP, SURVEYOR)
	return &SurveyorSocket{socket}, err
}

// Deadline returns the deadline for the surveyor. Default value is 1 second.
func (s *SurveyorSocket) Deadline() (time.Duration, error) {
	return s.Socket.SockOptDuration(C.NN_SURVEYOR, C.NN_SURVEYOR_DEADLINE, time.Millisecond)
}

// SetDeadline specifies how long to wait for responses to the survey. Once
// the deadline expires, receive function will return ETIMEDOUT error and
// all subsequent responses to the survey will be silently dropped.
func (s *SurveyorSocket) SetDeadline(deadline time.Duration) error {
	return s.Socket.SetSockOptDuration(C.NN_SURVEYOR, C.NN_SURVEYOR_DEADLINE, time.Millisecond, deadline)
}

type RespondentSocket struct {
	*Socket
}

// NewRespondentSocket creates a respondent socket used to respond to the
// survey. Survey is received using receive function, response is sent
// using send function. This socket can be connected to at most one peer.
func NewRespondentSocket() (*RespondentSocket, error) {
	socket, err := NewSocket(AF_SP, RESPONDENT)
	return &RespondentSocket{socket}, err
}
