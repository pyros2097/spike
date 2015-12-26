package main

import (
	"os"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	conn, err := websocket.Dial("ws://localhost:4000/ws", "", "http://localhost:4000")
	throwError("", err)
	for _ = range time.Tick(1 * time.Second) {
		err := websocket.JSON.Send(conn, "yello")
		throwError("Could not send message", err)
	}
	conn.Close()
	os.Exit(0)
}

func throwError(message string, err error) {
	if err != nil {
		panic(message+": %s", err)
	}
}
