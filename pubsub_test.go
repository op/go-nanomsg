// Go binding for nanomsg

package nanomsg

import (
	"testing"
	"time"
)

var socketAddress string = "inproc://a"

func TestPubSub(t *testing.T) {
	var err error
	var pub *PubSocket
	var sub1, sub2 *SubSocket

	if pub, err = NewPubSocket(); err != nil {
		t.Fatal(err)
	}
	if _, err = pub.Bind(socketAddress); err != nil {
		t.Fatal(err)
	}

	if sub1, err = NewSubSocket(); err != nil {
		t.Fatal(err)
	}
	if sub2, err = NewSubSocket(); err != nil {
		t.Fatal(err)
	}
	if err = sub1.Subscribe(""); err != nil {
		t.Fatal(err)
	}
	if err = sub2.Subscribe(""); err != nil {
		t.Fatal(err)
	}
	if _, err = sub1.Connect(socketAddress); err != nil {
		t.Fatal(err)
	}
	if _, err = sub2.Connect(socketAddress); err != nil {
		t.Fatal(err)
	}

	// Wait till connections are established to prevent message loss.
	time.Sleep(10 * time.Millisecond)

	nbytes, err := pub.Send([]byte("0123456789012345678901234567890123456789"), 0)
	if err != nil {
		t.Fatal(err)
	} else if nbytes != 40 {
		t.Fatal(nbytes)
	}

	buf, err := sub1.Recv(0)
	if err != nil {
		t.Fatal(err)
	} else if len(buf) != 40 {
		t.Fatal(len(buf))
	}

	buf, err = sub2.Recv(0)
	if err != nil {
		t.Fatal(err)
	} else if len(buf) != 40 {
		t.Fatal(len(buf))
	}

	if err = pub.Close(); err != nil {
		t.Fatal(err)
	}
	if err = sub1.Close(); err != nil {
		t.Fatal(err)
	}
	if err = sub2.Close(); err != nil {
		t.Fatal(err)
	}
}
