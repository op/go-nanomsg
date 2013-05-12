// Go binding for nanomsg

package nanomsg

// #include <nanomsg/nn.h>
// #include <nanomsg/pair.h>
// #include <nanomsg/reqrep.h>
// #include <stdlib.h>
// #cgo LDFLAGS: -lnanomsg
import "C"

import (
	"reflect"
	"runtime"
	"unsafe"
)

// C.NN_MSG is defined as size_t(-1), which makes cgo produce an error.
const nn_msg = ^C.size_t(0)

// SP address families.
type Domain int

const (
	SP     = Domain(C.AF_SP)
	SP_RAW = Domain(C.AF_SP_RAW)
)

type Protocol int

const (
	// reqrep.h
	REQ = Protocol(C.NN_REQ)
	REP = Protocol(C.NN_REP)

	// pair.h
	PAIR = Protocol(C.NN_PAIR)
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

func (s *Socket) Send(data []byte, flags int) error {
	buf := unsafe.Pointer(&data[0])
	length := C.size_t(len(data))
	if size, err := C.nn_send(s.socket, buf, length, C.int(flags)); size < 0 {
		return nnError(err)
	}

	return nil
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
			Len: capacity,
			Cap: capacity,
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

type Endpoint struct {
	Address string
	endpoint C.int
}

func (e *Endpoint) String() string {
	return e.Address
}
