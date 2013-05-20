// Go binding for nanomsg

package nanomsg

// #include <nanomsg/nn.h>
// #include <stdlib.h>
// #cgo LDFLAGS: -lnanomsg
import "C"

import (
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

// C.NN_MSG is defined as size_t(-1), which makes cgo produce an error.
const nn_msg = ^C.size_t(0)

// SP address families.
type Domain int

const (
	AF_SP     = Domain(C.AF_SP)
	AF_SP_RAW = Domain(C.AF_SP_RAW)
)

type Protocol int

// Sending and receiving can be controlled with these flags.
const (
	DontWait = int(C.NN_DONTWAIT)
)

type Socket struct {
	socket C.int
}

// Create a socket.
func NewSocket(domain Domain, protocol Protocol) (*Socket, error) {
	soc, err := C.nn_socket(C.int(domain), C.int(protocol))
	if err != nil {
		return nil, nnError(err)
	}

	socket := &Socket{soc}
	runtime.SetFinalizer(socket, (*Socket).Close)
	return socket, nil
}

// Close a socket.
func (s *Socket) Close() error {
	if s.socket != 0 {
		if rc, err := C.nn_close(s.socket); rc != 0 {
			return nnError(err)
		}
		s.socket = C.int(0)
	}
	return nil
}

// Add a local endpoint to the socket.
func (s *Socket) Bind(address string) (*Endpoint, error) {
	cstr := C.CString(address)
	defer C.free(unsafe.Pointer(cstr))
	eid, err := C.nn_bind(s.socket, cstr)
	if eid < 0 {
		return nil, nnError(err)
	}
	return &Endpoint{address, eid}, nil
}

// Add a remote endpoint to the socket.
func (s *Socket) Connect(address string) (*Endpoint, error) {
	cstr := C.CString(address)
	defer C.free(unsafe.Pointer(cstr))
	eid, err := C.nn_connect(s.socket, cstr)
	if eid < 0 {
		return nil, nnError(err)
	}
	return &Endpoint{address, eid}, nil
}

// Removes an endpoint from the socket.
func (s *Socket) Shutdown(endpoint *Endpoint) error {
	if rc, err := C.nn_shutdown(s.socket, endpoint.endpoint); rc != 0 {
		return nnError(err)
	}
	return nil
}

func (s *Socket) Send(data []byte, flags int) (int, error) {
	var buf unsafe.Pointer
	if len(data) != 0 {
		buf = unsafe.Pointer(&data[0])
	}
	length := C.size_t(len(data))
	size, err := C.nn_send(s.socket, buf, length, C.int(flags))
	if size < 0 {
		return int(size), nnError(err)
	}
	return int(size), nil
}

func (s *Socket) Recv(flags int) ([]byte, error) {
	var err error
	var buf unsafe.Pointer
	var length C.int

	if length, err = C.nn_recv(s.socket, unsafe.Pointer(&buf), nn_msg, C.int(flags)); length < 0 {
		return nil, nnError(err)
	}

	// TODO why is the latter variant faster than the zero copy variant?
	zeroCopy := true
	if zeroCopy {
		capacity := int(length)
		header := &reflect.SliceHeader{
			Data: uintptr(buf),
			Len:  capacity,
			Cap:  capacity,
		}
		data := *((*[]byte)(unsafe.Pointer(header)))
		runtime.SetFinalizer(&data, finalizeMsg)
		return data, nil
	} else {
		data := C.GoBytes(buf, length)
		if rc, err := C.nn_freemsg(buf); rc != 0 {
			return data, nnError(err)
		}
		return data, nil
	}
}

func finalizeMsg(datap *[]byte) error {
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(datap))
	if rc, err := C.nn_freemsg(unsafe.Pointer(hdrp.Data)); rc != 0 {
		return nnError(err)
	}
	return nil
}

func (s *Socket) GetSockOptInt(level, option C.int) (int, error) {
	var value C.int
	length := C.size_t(unsafe.Sizeof(value))
	rc, err := C.nn_getsockopt(s.socket, level, option, unsafe.Pointer(&value), &length)
	if rc != 0 {
		err = nnError(err)
		return int(value), err
	}
	return int(value), nil
}

func (s *Socket) SetSockOptInt(level, option C.int, value int) error {
	val := C.int(value)
	length := C.size_t(unsafe.Sizeof(val))
	rc, err := C.nn_setsockopt(s.socket, level, option, unsafe.Pointer(&val), length)
	if rc != 0 {
		return nnError(err)
	}
	return nil
}

// SetSockOptString sets the value of the option.
func (s *Socket) SetSockOptString(level, option C.int, value string) error {
	cstr := C.CString(value)
	defer C.free(unsafe.Pointer(cstr))
	length := C.size_t(len(value))
	rc, err := C.nn_setsockopt(s.socket, level, option, unsafe.Pointer(cstr), length)
	if rc != 0 {
		return nnError(err)
	}
	return nil
}

// Linger returns how long the socket should try to send pending outbound
// messages after Close() have been called, in nanoseconds (as defined by
// time.Duration). Negative value means infinite linger.
func (s *Socket) Linger() (time.Duration, error) {
	lingerMs, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_LINGER)
	linger := time.Duration(lingerMs) * time.Millisecond
	return linger, err
}

// SetLinger sets how long the socket should try to send pending outbound
// messages after Close() have been called, in nanoseconds (as defined by
// time.Duration). Negative value means infinite linger.
//
// Default value is 1 second.
func (s *Socket) SetLinger(linger time.Duration) error {
	lingerMs := int(linger / time.Millisecond)
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_LINGER, lingerMs)
}

// SendBuffer returns the size of the send buffer, in bytes. To
// prevent blocking for messages larger than the buffer, exactly one
// message may be buffered in addition to the data in the send buffer.
// Default value is 128kB.
func (s *Socket) SendBuffer() (int64, error) {
	size, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDBUF)
	return int64(size), err
}

// SetSendBuffer sets the send buffer size.
func (s *Socket) SetSendBuffer(sndBuf int64) error {
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDBUF, int(sndBuf))
}

// RecvBuffer returns the size of the receive buffer, in bytes. To
// prevent blocking for messages larger than the buffer, exactly one
// message may be buffered in addition to the data in the receive
// buffer. Default value is 128kB.
func (s *Socket) RecvBuffer() (int64, error) {
	size, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_RCVBUF)
	return int64(size), err
}

// SetRecvBuffer sets the receive buffer size.
func (s *Socket) SetRecvBuffer(rcvBuf int64) error {
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_RCVBUF, int(rcvBuf))
}

// SendTimeout returns the timeout for send operation on the socket.
// If message cannot be sent within the specified timeout, EAGAIN
// error is returned. Negative value means infinite timeout. Default
// value is -1.
func (s *Socket) SendTimeout() (time.Duration, error) {
	timeoutMs, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDTIMEO)
	timeout := time.Duration(timeoutMs)
	if timeout >= 0 {
		timeout *= time.Millisecond
	}
	return timeout, err
}

// SetSendTimeout sets the timeout for send operations.
func (s *Socket) SetSendTimeout(timeout time.Duration) error {
	timeoutMs := int(timeout / time.Millisecond)
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDTIMEO, timeoutMs)
}

// RecvTimeout returns the timeout for recv operation on the
// socket. If message cannot be received within the specified timeout,
// EAGAIN error is returned. Negative value means infinite timeout.
// Default value is -1.
func (s *Socket) RecvTimeout() (time.Duration, error) {
	timeoutMs, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_RCVTIMEO)
	timeout := time.Duration(timeoutMs)
	if timeout >= 0 {
		timeout *= time.Millisecond
	}
	return timeout, err
}

// SetRecvTimeout sets the timeout for recv operations.
func (s *Socket) SetRecvTimeout(timeout time.Duration) error {
	timeoutMs := int(timeout / time.Millisecond)
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_RCVTIMEO, timeoutMs)
}

// ReconnectInterval, for connection-based transports such as TCP,
// this option specifies how long to wait, when connection is broken
// before trying to re-establish it. Note that actual reconnect
// interval may be randomised to some extent to prevent severe
// reconnection storms. Default value is 0.1 second.
func (s *Socket) ReconnectInterval() (time.Duration, error) {
	ivlMs, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL)
	ivl := time.Duration(ivlMs) * time.Millisecond
	return ivl, err
}

// SetReconnectInterval sets the reconnect interval.
func (s *Socket) SetReconnectInterval(interval time.Duration) error {
	ivlMs := int(interval / time.Millisecond)
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL, ivlMs)
}

// ReconnectIntervalMax, together with ReconnectInterval, specifies
// maximum reconnection interval. On each reconnect attempt, the
// previous interval is doubled until this value is reached. Value of
// zero means that no exponential backoff is performed and reconnect
// interval is based only on the reconnect interval. If this value is
// less than the reconnect interval, it is ignored. Default value is
// 0.
func (s *Socket) ReconnectIntervalMax() (time.Duration, error) {
	ivlMs, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL_MAX)
	ivl := time.Duration(ivlMs) * time.Millisecond
	return ivl, err
}

// SetReconnectIntervalMax sets the maximum reconnect interval.
func (s *Socket) SetReconnectIntervalMax(interval time.Duration) error {
	ivlMs := int(interval / time.Millisecond)
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL_MAX, ivlMs)
}

// SendPrio sets outbound priority for endpoints subsequently added to
// the socket. This option has no effect on socket types that send
// messages to all the peers. However, if the socket type sends each
// message to a single peer (or a limited set of peers), peers with
// high priority take precedence over peers with low priority. The
// type of the option is int. Highest priority is 1, lowest priority
// is 16. Default value is 8.
func (s *Socket) SendPrio() (int, error) {
	return s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDPRIO)
}

// SetSendPrio sets the sending priority.
func (s *Socket) SetSendPrio(sndPrio int) error {
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDPRIO, sndPrio)
}

func (s *Socket) SendFd() (uintptr, error) {
	fd, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDFD)
	return uintptr(fd), err
}

func (s *Socket) RecvFd() (uintptr, error) {
	fd, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_RCVFD)
	return uintptr(fd), err
}

func (s *Socket) Domain() (Domain, error) {
	domain, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_DOMAIN)
	return Domain(domain), err
}

func (s *Socket) Protocol() (Protocol, error) {
	proto, err := s.GetSockOptInt(C.NN_SOL_SOCKET, C.NN_PROTOCOL)
	return Protocol(proto), err
}

type Endpoint struct {
	Address  string
	endpoint C.int
}

func (e *Endpoint) String() string {
	return e.Address
}
