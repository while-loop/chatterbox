package main

import (
	"log"
	"net"
	"strings"
	"time"
)

// connFunc is a function that will supply a map of current connected clients to
// hub. This function will be ran in and only in the goroutine created in Hub.Start().
type connFunc func(conns map[string]net.Conn)

type Hub struct {
	connChan chan connFunc
}

// NewHub creates a hub an initializes unexported vars
func NewHub() *Hub {
	return &Hub{
		connChan: make(chan connFunc),
	}
}

// OnMessage handles messages from Hub clients.
// When a message is received this function is called.
func (h *Hub) OnMessage(message string) {
	if strings.HasPrefix(message, "!time") {
		message = time.Now().String()
	}

	h.Broadcast(message)
}

// Broadcast sends message to all connected Hub clients.
func (h *Hub) Broadcast(message string) {
	h.connChan <- func(conns map[string]net.Conn) {
		for _, conn := range conns {
			go conn.Write([]byte(message))
		}
	}
}

// Register adds a user to to Hub.
// Here is where any auth or setup for a client will happen.
func (h *Hub) Register(conn net.Conn) {
	h.connChan <- func(conns map[string]net.Conn) {
		conns[conn.RemoteAddr().String()] = conn
		log.Printf("User %s joined hub! %d connected users", conn.RemoteAddr(), len(conns))
		conn.Write([]byte("welcome!"))
	}
}

// Deregister handles any client cleanup.
// Note, Deregister does not close the underlying client connection,
// that is handled by the connection lifecycle implementation.
func (h *Hub) Deregister(conn net.Conn) {
	h.connChan <- func(conns map[string]net.Conn) {
		delete(conns, conn.RemoteAddr().String())
		log.Printf("User %s left hub! %d connected users", conn.RemoteAddr(), len(conns))
	}
}

// Start creates a map of connections and waits for functions to be passed
// by the connChan var.
func (h *Hub) Start() {
	conns := map[string]net.Conn{}

	for fn := range h.connChan {
		fn(conns)
	}
}
