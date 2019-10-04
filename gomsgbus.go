package main

import (
	"fmt"
	"os"
	"time"

	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/bus"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func node(args []string) {
	var sock mangos.Socket
	var err error
	var msg []byte
	var x int

	if sock, err = bus.NewSocket(); err != nil {
		die("bus.NewSocket: %s", err)
	}
	if err = sock.Listen(args[2]); err != nil {
		die("sock.Listen: %s", err.Error())
	}

	// wait for everyone to start listening
	time.Sleep(time.Second)
	for x = 3; x < len(args); x++ {
		if err = sock.Dial(args[x]); err != nil {
			die("socket.Dial: %s", err.Error())
		}
	}

	// wait for everyone to join
	time.Sleep(time.Second)

	fmt.Printf("%s: SENDING '%s' ONTO BUS\n", args[1], args[1])
	if err = sock.Send([]byte(args[1])); err != nil {
		die("sock.Send: %s", err.Error())
	}
	for {
		if msg, err = sock.Recv(); err != nil {
			die("sock.Recv: %s", err.Error())
		}
		fmt.Printf("%s: RECEIVED \"%s\" FROM BUS\n", args[1],
			string(msg))

	}
}

func main() {
	if len(os.Args) > 3 {
		node(os.Args)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Usage: bus <NODENAME> <URL> <URL>... \n")
	os.Exit(1)
}
