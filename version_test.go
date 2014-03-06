// Go binding for nanomsg

package nanomsg

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// nanomsg-0.1 C:R:A(0:0:0)
	// nanomsg-0.2 C:R:A(0:1:0)
	// nanomsg-0.3 C:R:A(1:0:1)
	if Version.Current < 0 || Version.Revision < 0 || Version.Age < 0 {
		t.Fatalf("Unexpected library version: %s", Version)
	}
	// if change current, set revision=0
	if Version.Current > 0 {
		if Version.Revision != 0 {
			t.Fatalf("Unexpected library version: %s", Version)
		}
	}
}
