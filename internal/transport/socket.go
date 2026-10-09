package transport

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
)

type RegisterConnection func() string

type UnregisterConnection func(connectionName string)

type ProcessMessage func(connectionName string, line string)

var (
	connections   = make(map[string]net.Conn)
	connectionsMu sync.RWMutex
)

func ServeSocket(
	socketPath string,
	register RegisterConnection,
	unregister UnregisterConnection,
	process ProcessMessage,
) error {
	if err := os.RemoveAll(socketPath); err != nil {
		return err
	}
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	defer func() { _ = listener.Close() }()
	log.Println("Listening on", socketPath)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handleConnection(conn, register, unregister, process)
	}
}

func handleConnection(
	conn net.Conn,
	register RegisterConnection,
	unregister UnregisterConnection,
	process ProcessMessage,
) {
	name := register()
	connectionsMu.Lock()
	connections[name] = conn
	connectionsMu.Unlock()
	defer func() {
		connectionsMu.Lock()
		delete(connections, name)
		connectionsMu.Unlock()
		unregister(name)
		_ = conn.Close()
	}()
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 4096), 1<<20)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("received: %q\n", message)
		process(name, message)
	}
	if err := scanner.Err(); err != nil {
		log.Println("read:", err)
	}
}

func SendOverSocket(connectionName string, message string) {
	connectionsMu.RLock()
	conn, ok := connections[connectionName]
	connectionsMu.RUnlock()
	if !ok {
		log.Println("connection not found:", connectionName)
		return
	}
	if _, err := conn.Write([]byte(message + "\n")); err != nil {
		log.Println("write:", err)
	}
}
