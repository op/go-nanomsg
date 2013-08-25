package nanomsg

// #include <nanomsg/nn.h>
import "C"

// Version holds the nanomsg version which is used. nanomsg uses libtool's
// versioning system.
var Version = struct {
	Current  int
	Revision int
	Age      int
}{
	int(C.NN_VERSION_CURRENT),
	int(C.NN_VERSION_REVISION),
	int(C.NN_VERSION_AGE),
}
