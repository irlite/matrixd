package main

import (
	"log"

	"github.com/irlite/matrixd/internal/hub"
	"github.com/irlite/matrixd/internal/matrix/session"
	"github.com/irlite/matrixd/internal/transport"
)

func main() {
	session.SetNotify(hub.SendMessage)
	hub.RegisterSender(hub.TypeSocket, transport.SendOverSocket)
	err := transport.ServeSocket(
		"/tmp/matrixd.sock",
		func() string {
			return hub.RegisterConnection(hub.TypeSocket)
		},
		hub.UnregisterConnection,
		hub.ProcessIncomingMessage,
	)
	if err != nil {
		log.Fatal(err)
	}
}
