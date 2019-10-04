package main

import (
	"fmt"
	"log"
	"nanomsg.org/go/mangos/v2"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"nanomsg.org/go/mangos/v2/protocol/bus"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

var (
	nodename            = os.Args[1]
	cmdport             = os.Args[2]
	myurl               = os.Args[3]
	firstPeerIdx        = 4
	minArgLength        = 5
	incomingBusChan     = make(chan string)
	incomingHttpCmdChan = make(chan string)
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
	go msgHub(busSocket)
	go busReceiver(busSocket)

	e := httpSetup()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":" + cmdport)))

}

func busReceiver(busSocket mangos.Socket) {
	// Start receiving...
	for {
		msg, err := busSocket.Recv()
		if err != nil {
			die("sock.Recv: %s", err.Error())
		}
		incomingBusChan <- string(msg)
	}
}

func msgHub(busSocket mangos.Socket) {

	for {
		select {
		case httpMsg := <-incomingHttpCmdChan:
			err := busSocket.Send([]byte(httpMsg))
			if err != nil {
				fmt.Errorf("cannot send http cmd to bus: %#v", err)
			}

		case busMsg := <-incomingBusChan:
			fmt.Printf("%s: RECEIVED \"%s\" FROM BUS\n", nodename, busMsg)
		}
	}

}

func httpSetup() *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello GoMsgBus!")
	})

	e.POST("/cmd", func(c echo.Context) error {
		incomingHttpCmdChan <- fmt.Sprintf("The time now is: %v", time.Now())
		return c.String(http.StatusOK, "Kommando erfolgreich gepostet...")
	})
	e.HideBanner = true
	return e
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

	log.Println(nodename+" is waiting for peers....")
	time.Sleep(time.Second)
	log.Println("Continue....")
	// connect to peers
	for x := firstPeerIdx; x < len(os.Args); x++ {
		if err = busSocket.Dial(os.Args[x]); err != nil {
			die("socket.Dial: %s", err.Error())
		}
	}

	return busSocket

}
