package transport

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
)

type Handler func(connName string, line string) (response string, newConnName string)

var (
	connections   = make(map[string]net.Conn)
	connectionsMu sync.RWMutex
)

func ServeSocket(socketPath string, handle Handler) error {
	if err := os.RemoveAll(socketPath); err != nil {
		return err
	}
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("Listening on", socketPath)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handleConnection(conn, handle)
	}
}

func handleConnection(conn net.Conn, handle Handler) {
	var connName string
	defer func() {
		if connName != "" {
			connectionsMu.Lock()
			if connections[connName] == conn {
				delete(connections, connName)
			}
			connectionsMu.Unlock()
		}
		conn.Close()
	}()
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 4096), 1<<20)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("received: %q\n", line)
		response, newConnName := handle(connName, line)
		if newConnName != "" {
			connName = newConnName
			connectionsMu.Lock()
			connections[connName] = conn
			connectionsMu.Unlock()
		}
		if _, err := conn.Write([]byte(response + "\n")); err != nil {
			log.Println("write:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("read:", err)
	}
}

func SendOverSocket(connName string, message string) {
	connectionsMu.RLock()
	conn, ok := connections[connName]
	connectionsMu.RUnlock()
	if !ok {
		log.Println("connection not found:", connName)
		return
	}
	if _, err := conn.Write([]byte(message + "\n")); err != nil {
		log.Println("write:", err)
	}
}
