// Go binding for nanomsg

package nanomsg

// #include <nanomsg/nn.h>
import "C"

import (
	"errors"
	"syscall"
)

// TODO make sure EINTR is returned as expected
func nnError(err error) error {
	if errno, ok := err.(syscall.Errno); ok {
		return errors.New(C.GoString(C.nn_strerror(C.int(errno))))
	}
	return err
}
