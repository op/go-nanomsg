## Golang nanomsg bindings

Package nanomsg adds language bindings for nanomsg in Go. nanomsg is a
high-performance implementation of several "scalability protocols". See
http://nanomsg.org/ for more information.

This is a work in progress. nanomsg is still in a beta stage. Expect its
API, or this binding, to change.

## Installing

This is a cgo based library and requires the nanomsg library to build. Install
it either from [source](http://nanomsg.org/download.html) or use your package
manager of choice. 0.9 or later is required.

### Using *go get*

    $ go get github.com/op/go-nanomsg

After this command *go-nanomsg* is ready to use. Its source will be in:

    $GOROOT/src/pkg/github.com/op/go-nanomsg

You can use `go get -u -a` to update all installed packages.

## Documentation

For docs, see http://godoc.org/github.com/op/go-nanomsg or run:

    $ go doc github.com/op/go-nanomsg

## Alternatives

There is now also an implementation of nanomsg in pure Go. See
https://github.com/gdamore/mangos for more details.
