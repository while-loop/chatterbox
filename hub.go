package main

import (
	"log"
	"net"
	"strings"
	"time"
)

type connFunc func(conns map[string]net.Conn)

type Hub struct {
	conns    map[string]net.Conn
	connChan chan connFunc
}

func NewHub() *Hub {
	return &Hub{
		conns:    map[string]net.Conn{},
		connChan: make(chan connFunc),
	}
}

func (h *Hub) OnMessage(message string) {
	if strings.HasPrefix(message, "!time") {
		message = time.Now().String()
	}

	h.Broadcast(message)
}

func (h *Hub) Broadcast(message string) {
	h.connChan <- func(conns map[string]net.Conn) {
		for _, conn := range conns {
			go conn.Write([]byte(message))
		}
	}
}

func (h *Hub) Register(conn net.Conn) {
	h.connChan <- func(conns map[string]net.Conn) {
		conns[conn.RemoteAddr().String()] = conn
		log.Printf("User %s joined hub! %d connected users", conn.RemoteAddr(), len(conns))
		conn.Write([]byte("welcome!"))
	}
}

func (h *Hub) Deregister(conn net.Conn) {
	h.connChan <- func(conns map[string]net.Conn) {
		delete(conns, conn.RemoteAddr().String())
		log.Printf("User %s left hub! %d connected users", conn.RemoteAddr(), len(conns))
	}
}

func (h *Hub) Start() error {
	for fn := range h.connChan {
		fn(h.conns)
	}
	return nil
}
