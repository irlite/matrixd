package logic

import (
	"fmt"
	"sync"
)

type IPCType string

const (
	TypeSocket IPCType = "socket"
	TypeHTTP   IPCType = "http"
)

type Sender func(connName string, message string)

type MessageType string

const (
	MatrixdDBData   MessageType = "matrixdDBData"
	MessageRegister MessageType = "register"
	MessageSend     MessageType = "SEND"
	MessagePing     MessageType = "PING"
)

var (
	mu                  sync.RWMutex
	senders             = make(map[IPCType]Sender)
	connectionToIPCType = make(map[string]IPCType)
)

func RegisterSender(t IPCType, s Sender) {
	mu.Lock()
	senders[t] = s
	mu.Unlock()
}

func ProcessIncomingMessage(
	connName string,
	message string,
	t IPCType,
) (string, string) {
	parsedMessage := parseMessage(message)
	switch parsedMessage.messageType {
	case MessageRegister:
		fmt.Println("register")
	case MessageSend:
		fmt.Println("send")
	case MessagePing:
		fmt.Println("ping")
		return "PONG", ""
	default:
		fmt.Println("unknown command")
	}
	return "", ""
}

func SendMessage(connName string, message string) {
	mu.RLock()
	send := senders[connectionToIPCType[connName]]
	mu.RUnlock()
	if send != nil {
		send(connName, message)
	}
}
