## Golang nanomsg bindings

Package nanomsg adds language bindings for nanomsg in Go. nanomsg is a
high-performance implementation of several "scalability protocols". See
http://nanomsg.org/ for more information.

This is a work in progress. nanomsg is still in a beta stage. Expect its
API, or this binding, to change.

## Installing

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
https://bitbucket.org/gdamore/mangos for more details.
