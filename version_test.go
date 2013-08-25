// Go binding for nanomsg

package nanomsg

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version.Current != 0 {
		t.Fatalf("Unexpected library version: %s", Version)
	}
}
