package main

import (
	"bytes"
	"flag"
	"os"
	"strconv"

	"../.."
)

func main() {
	flag.Parse()

	if flag.NArg() != 3 {
		println("usage: remote_thr <connect-to> <msg-size> <msg-count")
		os.Exit(1)
	}

	connectTo := flag.Arg(0)
	sz, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		panic(err)
	}
	count, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		panic(err)
	}

	// TODO define pair socket?
	s, err := nanomsg.NewSocket(nanomsg.SP, nanomsg.PAIR)
	if err != nil {
		panic(err)
	}
	_, err = s.Connect(connectTo)
	if err != nil {
		panic(err)
	}

	buf := bytes.Repeat([]byte{111}, sz)
	nbytes, err := s.Send(buf[0:0], 0)
	if err != nil {
		panic(err)
	} else if nbytes != 0 {
		panic(nbytes)
	}

	for i := 0; i < count; i++ {
		nbytes, err = s.Send(buf, 0)
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
