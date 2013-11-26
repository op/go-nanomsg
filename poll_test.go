// Go binding for nanomsg

package nanomsg

import (
	"testing"
	"time"
)

func TestPoller(t *testing.T) {
	sa, err := NewPairSocket()
	if err != nil {
		t.Fatal(err)
	}
	defer sa.Close()
	if _, err := sa.Bind("inproc://a"); err != nil {
		t.Fatal(err)
	}
	sb, err := NewPairSocket()
	if err != nil {
		t.Fatal(err)
	}
	defer sb.Close()
	if _, err := sb.Connect("inproc://a"); err != nil {
		t.Fatal(err)
	}

	// Create a poller which makes sure it can send on both sockets.
	var poller Poller
	pia := poller.Add(sa.Socket, true, true)
	pib := poller.Add(sb.Socket, true, true)

	if i, err := poller.Poll(10 * time.Millisecond); err != nil {
		t.Fatal(err)
	} else if i != 2 {
		t.Error("should be able to send", i)
	}

	// Remove the send check and pass in a message on one of the sockets.
	pia.PollSend(false)
	pib.PollSend(false)
	if _, err := sa.Send([]byte("abc"), DontWait); err != nil {
		t.Fatal(err)
	}

	if i, err := poller.Poll(10 * time.Millisecond); err != nil {
		t.Fatal(err)
	} else if i != 1 {
		t.Error("no events", i)
	}

	if pia.CanRecv() {
		t.Error("unexpected to recv")
	}
	if !pib.CanRecv() {
		t.Error("expected recv")
	}

	// Receieve the message and allow the poller to timeout.
	if _, err := sb.Recv(DontWait); err != nil {
		t.Fatal(err)
	}
	if i, err := poller.Poll(1 * time.Millisecond); err != nil {
		t.Fatal(err)
	} else if i != 0 {
		t.Error("no events should be available")
	}
}
