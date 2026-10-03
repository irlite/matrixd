package main

import (
	"log"
	"matrixd/internal/logic"
	"matrixd/internal/transport"
)

func main() {
	logic.RegisterSender(logic.TypeSocket, transport.SendOverSocket)
	err := transport.ServeSocket("/tmp/matrixd.sock", func(connName, line string) (string, string) {
		var response, newConnName string = logic.ProcessIncomingMessage(connName, line, logic.TypeSocket)
		return response, newConnName
	})
	if err != nil {
		log.Fatal(err)
	}
}
