// Go binding for nanomsg

package nanomsg

// #include <nanomsg/nn.h>
import "C"

import (
	"syscall"
)

// Errno defines specific nanomsg errors
//
// The errors returned from operations on the nanomsg library and the Go
// bindings for it tries to return all errors using the errors already found in
// Go like syscall.EADDRINUSE. There are some errors that only exists in
// nanomsg and these are defined as Errno.
type Errno syscall.Errno

const (
	nn_hausnumero = Errno(int(C.NN_HAUSNUMERO))

	// Nanomsg specific errors
	ETERM = Errno(int(C.ETERM))
	EFSM  = Errno(int(C.EFSM))
)

var errorStrings = map[Errno]string{
	ETERM: C.GoString(C.nn_strerror(C.ETERM)),
	EFSM:  C.GoString(C.nn_strerror(C.EFSM)),
}

func (e Errno) Error() string {
	s := errorStrings[e]
	if s == "" {
		s = C.GoString(C.nn_strerror(C.int(e)))
	}
	return s
}

// Errors expected to be found when calling the nanomsg library should map
// directly to the syscall.Errno error found in Go. On some platforms, the
// POSIX error is not defined and will most likely be different from the one
// define in the nanomsg library as in the Go world. This map is used to
// automatically make errors found in the nanomsg library to map nicely to the
// ones users would expect in Go.
var errnoMap = map[syscall.Errno]syscall.Errno{
	syscall.Errno(int(C.ENOTSUP)): syscall.ENOTSUP,
	syscall.Errno(int(C.EPROTONOSUPPORT)): syscall.EPROTONOSUPPORT,
	syscall.Errno(int(C.ENOBUFS)): syscall.ENOBUFS,
	syscall.Errno(int(C.ENETDOWN)): syscall.ENETDOWN,
	syscall.Errno(int(C.EADDRINUSE)): syscall.EADDRINUSE,
	syscall.Errno(int(C.EADDRNOTAVAIL)): syscall.EADDRNOTAVAIL,
	syscall.Errno(int(C.ECONNREFUSED)): syscall.ECONNREFUSED,
	syscall.Errno(int(C.EINPROGRESS)): syscall.EINPROGRESS,
	syscall.Errno(int(C.ENOTSOCK)): syscall.ENOTSOCK,
	syscall.Errno(int(C.EAFNOSUPPORT)): syscall.EAFNOSUPPORT,
	syscall.Errno(int(C.EPROTO)): syscall.EPROTO,
}

// nnError takes an error returned from the C nanomsg library and transforms it
// to standard errors when possible.
func nnError(err error) error {
	if errno, ok := err.(syscall.Errno); ok {
		// nn_hausnumero is the number the nanomsg library has chosen to hopefully be
		// different enough to not collide with any other errno. try to convert it to
		// a known syscall error if found. If not found, it should hopefully be one
		// of the ones defined as nanomsg.Errno.
		if int(errno) >= int(nn_hausnumero) {
			sysErr, present := errnoMap[errno]
			if present {
				err = sysErr
			} else {
				return Errno(errno)
			}
		}
	}
	return err
}
