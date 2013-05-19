// Go binding for nanomsg

package nanomsg

import (
	"syscall"
	"testing"
)

func TestError(t *testing.T) {
	var err error
	var s *Socket
	if s, err = NewSocket(SP, REP); err != nil {
		t.Fatal(err)
	}
	// TODO inproc://a isn't available, even tho close is called?
	if _, err = s.Bind("inproc://b"); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Bind("inproc://b"); err == nil {
		t.Fatal("expected failure")
	} else {
		if err != syscall.EADDRINUSE {
			t.Fatal(err)
		}
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
}
