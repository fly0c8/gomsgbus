package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/bus"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

var (
	nodename     = os.Args[1]
	cmdport      = os.Args[2]
	myurl        = os.Args[3]
	firstPeerIdx = 4
	minArgLength = 5
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {

	if len(os.Args) < minArgLength {
		fmt.Fprintf(os.Stderr, "Usage: gomsgbus <NODENAME> <CMDPORT> <MYURL> <URL>... \n")
		os.Exit(1)
	}

	busSocket := busSetup()
	go busHandler(busSocket)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello GoMsgBus!")
	})
	e.POST("/cmd", func(c echo.Context) error {
		err := busSocket.Send([]byte(fmt.Sprintf("The time now is: %v", time.Now())))
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%#v", err))
		}
		return c.String(http.StatusOK, "Kommando erfolgreich gepostet...")
	})
	e.HideBanner = true
	e.Logger.Fatal(e.Start(fmt.Sprintf(":" + cmdport)))

}
func busHandler(busSocket mangos.Socket) {
	// Start receiving...
	for {
		msg, err := busSocket.Recv()
		if err != nil {
			die("sock.Recv: %s", err.Error())
		}
		fmt.Printf("%s: RECEIVED \"%s\" FROM BUS\n", nodename, string(msg))
	}
}
func busSetup() mangos.Socket {
	// start bus listener
	var busSocket mangos.Socket
	var err error

	if busSocket, err = bus.NewSocket(); err != nil {
		die("bus.NewSocket: %s", err)
	}

	if err = busSocket.Listen(myurl); err != nil {
		die("busSocket.Listen: %s", err.Error())
	}

	// wait for everyone to start listening
	time.Sleep(time.Second)

	// connect to peers
	for x := firstPeerIdx; x < len(os.Args); x++ {
		if err = busSocket.Dial(os.Args[x]); err != nil {
			die("socket.Dial: %s", err.Error())
		}
	}

	return busSocket

}
