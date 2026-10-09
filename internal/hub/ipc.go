package hub

import (
	"strconv"
	"sync"

	"github.com/irlite/matrixd/internal/wire"
)

type IPCType string

const (
	TypeSocket IPCType = "socket"
	TypeHTTP   IPCType = "http"
)

type Sender func(connName string, message string)

var (
	mu                  sync.RWMutex
	senders             = make(map[IPCType]Sender)
	connectionToIPCType = make(map[string]IPCType)
	nextConnectionID    int
)

func RegisterSender(ipc IPCType, sender Sender) {
	mu.Lock()
	senders[ipc] = sender
	mu.Unlock()
}

func RegisterConnection(ipc IPCType) string {
	mu.Lock()
	defer mu.Unlock()
	name := strconv.Itoa(nextConnectionID)
	nextConnectionID++
	connectionToIPCType[name] = ipc
	return name
}

func UnregisterConnection(connName string) {
	mu.Lock()
	delete(connectionToIPCType, connName)
	mu.Unlock()
}

func ProcessIncomingMessage(connName string, message string) {
	req, err := parseMessage(message)
	req.ConnName = connName
	if err != nil {
		SendMessage(connName, fail(req, err.Error()))
		return
	}
	SendMessage(connName, routeMessage(req))
}

func SendMessage(connName string, out wire.Outgoing) {
	mu.RLock()
	send := senders[connectionToIPCType[connName]]
	mu.RUnlock()
	if send == nil {
		return
	}
	send(connName, marshal(out))
}
