package main

import (
	"log"
	"net/http"

	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type connFunc func(conns map[string]*websocket.Conn)

var (
	wsUpgrader = websocket.Upgrader{}
)

type Hub struct {
	conns    map[string]*websocket.Conn
	connChan chan connFunc
}

func NewHub() *Hub {
	return &Hub{
		conns:    map[string]*websocket.Conn{},
		connChan: make(chan connFunc),
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("New request from", r.RemoteAddr)
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(r.RemoteAddr, err)
		return
	}

	h.Register(conn)
	for {
		kind, bs, err := conn.ReadMessage()
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			break
		}

		if kind == websocket.TextMessage {
			h.OnMessage(string(bs))
		}
	}
	h.Deregister(conn)
}

func (h *Hub) OnMessage(message string) {
	if strings.HasPrefix(message, "!time") {
		message = time.Now().String()
	}

	h.connChan <- func(conns map[string]*websocket.Conn) {
		for _, conn := range conns {
			go conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.connChan <- func(conns map[string]*websocket.Conn) {
		log.Printf("User %s joined hub! %d connected users", conn.RemoteAddr(), len(conns)+1)
		conns[conn.RemoteAddr().String()] = conn
		conn.WriteMessage(websocket.TextMessage, []byte("welcome!"))
	}
}

func (h *Hub) Deregister(conn *websocket.Conn) {
	h.connChan <- func(conns map[string]*websocket.Conn) {
		log.Printf("User %s left hub! %d connected users", conn.RemoteAddr(), len(conns)-1)
		delete(conns, conn.RemoteAddr().String())
	}
}

func (h *Hub) Start() error {
	for fn := range h.connChan {
		fn(h.conns)
	}
	return nil
}
