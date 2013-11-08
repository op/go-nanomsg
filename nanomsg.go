// Go binding for nanomsg

package nanomsg

// #include <nanomsg/nn.h>
// #include <stdlib.h>
// #cgo LDFLAGS: -lnanomsg
import "C"

import (
	"reflect"
	"runtime"
	"syscall"
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
	// socket is the actual nanomsg C API object
	socket C.int
}

// Create a socket.
func NewSocket(domain Domain, protocol Protocol) (*Socket, error) {
	soc, err := C.nn_socket(C.int(domain), C.int(protocol))
	if soc == -1 {
		return nil, nnError(err)
	}

	// Create the socket object and make sure we call Close before freeing up the
	// memory inside the Go runtime.
	socket := &Socket{socket: soc}
	socket.setFinalizer()
	return socket, nil
}

func (s *Socket) setFinalizer() {
	runtime.SetFinalizer(s, (*Socket).Close)
}

// Close a socket.
func (s *Socket) Close() error {
	if rc, err := C.nn_close(s.socket); rc != 0 {
		// If the close call was interrupted by the signal handler, nanomsg
		// would return EINTR. All is good except when Close() is called by the
		// Go runtime during garbage collection. When this happens, the Go
		// runtime clears the finalizer before running it. Hence we need to set
		// it again to avoid leaking resources.
		//
		// This has a couple of side-effects that might not be obvious.
		//
		// * If the user calls Close() two times, and the latter one is
		//   interrupted we will queue yet another call to Close. (Most likely
		//   the second call to nn_close will not block.)
		//
		// * If the user has set a custom finalizer, we would at this point
		//   override it.
		//
		// However, all of these scenarios is an unexpected use of this library.
		if err = nnError(err); err == syscall.EINTR {
			s.setFinalizer()
		}
		return err
	}
	// Once the socket has been closed, we no longer need to call Close when the
	// object is garbage collected.
	runtime.SetFinalizer(s, nil)
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

func (s *Socket) SockOptInt(level, option C.int) (int, error) {
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

// SockOptDuration retrieves the socket option as duration. unit is
// used to specify the unit which nanomsg exposes the option as.
func (s *Socket) SockOptDuration(level, option C.int, unit time.Duration) (time.Duration, error) {
	value, err := s.SockOptInt(level, option)
	return time.Duration(value) * unit, err
}

// SetSockOptDuration sets the socket option as duration. unit is
// used to specify the unit which nanomsg exposes the option as.
func (s *Socket) SetSockOptDuration(level, option C.int, unit, value time.Duration) error {
	return s.SetSockOptInt(level, option, int(value/unit))
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
// messages after Close() have been called. Negative value means
// infinite linger.
func (s *Socket) Linger() (time.Duration, error) {
	return s.SockOptDuration(C.NN_SOL_SOCKET, C.NN_LINGER, time.Millisecond)
}

// SetLinger sets how long the socket should try to send pending outbound
// messages after Close() have been called, in nanoseconds (as defined by
// time.Duration). Negative value means infinite linger.
//
// Default value is 1 second.
func (s *Socket) SetLinger(linger time.Duration) error {
	return s.SetSockOptDuration(C.NN_SOL_SOCKET, C.NN_LINGER, time.Millisecond, linger)
}

// SendBuffer returns the size of the send buffer, in bytes. To
// prevent blocking for messages larger than the buffer, exactly one
// message may be buffered in addition to the data in the send buffer.
// Default value is 128kB.
func (s *Socket) SendBuffer() (int64, error) {
	size, err := s.SockOptInt(C.NN_SOL_SOCKET, C.NN_SNDBUF)
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
	size, err := s.SockOptInt(C.NN_SOL_SOCKET, C.NN_RCVBUF)
	return int64(size), err
}

// SetRecvBuffer sets the receive buffer size.
func (s *Socket) SetRecvBuffer(rcvBuf int64) error {
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_RCVBUF, int(rcvBuf))
}

// SendTimeout returns the timeout for send operation on the socket.
// If message cannot be sent within the specified timeout, EAGAIN
// error is returned. Negative value means infinite timeout. Default
// value is infinite.
func (s *Socket) SendTimeout() (time.Duration, error) {
	return s.SockOptDuration(C.NN_SOL_SOCKET, C.NN_SNDTIMEO, time.Millisecond)
}

// SetSendTimeout sets the timeout for send operations.
func (s *Socket) SetSendTimeout(timeout time.Duration) error {
	return s.SetSockOptDuration(C.NN_SOL_SOCKET, C.NN_SNDTIMEO, time.Millisecond, timeout)
}

// RecvTimeout returns the timeout for recv operation on the
// socket. If message cannot be received within the specified timeout,
// EAGAIN error is returned. Negative value means infinite timeout.
// Default value is infinite.
func (s *Socket) RecvTimeout() (time.Duration, error) {
	return s.SockOptDuration(C.NN_SOL_SOCKET, C.NN_RCVTIMEO, time.Millisecond)
}

// SetRecvTimeout sets the timeout for recv operations.
func (s *Socket) SetRecvTimeout(timeout time.Duration) error {
	return s.SetSockOptDuration(C.NN_SOL_SOCKET, C.NN_RCVTIMEO, time.Millisecond, timeout)
}

// ReconnectInterval, for connection-based transports such as TCP,
// this option specifies how long to wait, when connection is broken
// before trying to re-establish it. Note that actual reconnect
// interval may be randomised to some extent to prevent severe
// reconnection storms. Default value is 0.1 second.
func (s *Socket) ReconnectInterval() (time.Duration, error) {
	return s.SockOptDuration(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL, time.Millisecond)
}

// SetReconnectInterval sets the reconnect interval.
func (s *Socket) SetReconnectInterval(interval time.Duration) error {
	return s.SetSockOptDuration(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL, time.Millisecond, interval)
}

// ReconnectIntervalMax, together with ReconnectInterval, specifies
// maximum reconnection interval. On each reconnect attempt, the
// previous interval is doubled until this value is reached. Value of
// zero means that no exponential backoff is performed and reconnect
// interval is based only on the reconnect interval. If this value is
// less than the reconnect interval, it is ignored. Default value is
// 0.
func (s *Socket) ReconnectIntervalMax() (time.Duration, error) {
	return s.SockOptDuration(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL_MAX, time.Millisecond)
}

// SetReconnectIntervalMax sets the maximum reconnect interval.
func (s *Socket) SetReconnectIntervalMax(interval time.Duration) error {
	return s.SetSockOptDuration(C.NN_SOL_SOCKET, C.NN_RECONNECT_IVL_MAX, time.Millisecond, interval)
}

// SendPrio sets outbound priority for endpoints subsequently added to
// the socket. This option has no effect on socket types that send
// messages to all the peers. However, if the socket type sends each
// message to a single peer (or a limited set of peers), peers with
// high priority take precedence over peers with low priority. The
// type of the option is int. Highest priority is 1, lowest priority
// is 16. Default value is 8.
func (s *Socket) SendPrio() (int, error) {
	return s.SockOptInt(C.NN_SOL_SOCKET, C.NN_SNDPRIO)
}

// SetSendPrio sets the sending priority.
func (s *Socket) SetSendPrio(sndPrio int) error {
	return s.SetSockOptInt(C.NN_SOL_SOCKET, C.NN_SNDPRIO, sndPrio)
}

func (s *Socket) SendFd() (uintptr, error) {
	fd, err := s.SockOptInt(C.NN_SOL_SOCKET, C.NN_SNDFD)
	return uintptr(fd), err
}

func (s *Socket) RecvFd() (uintptr, error) {
	fd, err := s.SockOptInt(C.NN_SOL_SOCKET, C.NN_RCVFD)
	return uintptr(fd), err
}

func (s *Socket) Domain() (Domain, error) {
	domain, err := s.SockOptInt(C.NN_SOL_SOCKET, C.NN_DOMAIN)
	return Domain(domain), err
}

func (s *Socket) Protocol() (Protocol, error) {
	proto, err := s.SockOptInt(C.NN_SOL_SOCKET, C.NN_PROTOCOL)
	return Protocol(proto), err
}

type Endpoint struct {
	Address  string
	endpoint C.int
}

func (e *Endpoint) String() string {
	return e.Address
}
