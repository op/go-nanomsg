package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/op/go-nanomsg"
)

func main() {
	flag.Parse()

	if flag.NArg() != 3 {
		println("usage: local_lat <bind-to> <msg-size> <roundtrips>")
		os.Exit(1)
	}

	bindTo := flag.Arg(0)
	sz, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		panic(err)
	}
	rts, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		panic(err)
	}

	s, err := nanomsg.NewPairSocket()
	if err != nil {
		panic(err)
	}
	if err = s.SetTCPNoDelay(true); err != nil {
		panic(err)
	}
	_, err = s.Bind(bindTo)
	if err != nil {
		panic(err)
	}

	for i := 0; i < rts; i++ {
		buf, err := s.Recv(0)
		if err != nil {
			panic(err)
		} else if len(buf) != sz {
			panic(sz)
		}

		nbytes, err := s.Send(buf, 0)
		if err != nil {
			panic(err)
		} else if nbytes != sz {
			panic(nbytes)
		}
	}

	err = s.Close()
	if err != nil {
		panic(err)
	}
}
