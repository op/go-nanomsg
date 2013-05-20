// Go binding for nanomsg

package nanomsg

// #include <nanomsg/nn.h>
import "C"

var symbols = make(map[string]int)

// Symbol returns the integer value for the given symbol.
func symbol(name string) (int, bool) {
	value, exists := symbols[name]
	return value, exists
}

func init() {
	var value C.int
	for i := 0; ; i++ {
		name, _ := C.nn_symbol(C.int(i), &value)
		if name == nil {
			break
		}
		symbols[C.GoString(name)] = int(value)
	}
}
