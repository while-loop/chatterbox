package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	wsUpgrader = websocket.Upgrader{}
)

type WSConn struct {
	conn   *websocket.Conn
	sendMu sync.Mutex
}

func NewWSConn(conn *websocket.Conn) *WSConn {
	return &WSConn{conn: conn, sendMu: sync.Mutex{}}
}

func FromHTTP(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(r.RemoteAddr, err)
			return
		}
		wsconn := NewWSConn(conn)

		defer func() {
			hub.Deregister(wsconn)
			wsconn.conn.Close()
		}()

		hub.Register(wsconn)
		for {
			kind, bs, err := wsconn.conn.ReadMessage()
			if err != nil {
				log.Println(wsconn.conn.RemoteAddr(), err)
				break
			}

			if kind == websocket.TextMessage {
				hub.OnMessage(string(bs))
			}
		}
	}
}

func (w *WSConn) Read(b []byte) (n int, err error) {
	panic("not used")
}

func (w *WSConn) Write(b []byte) (n int, err error) {
	w.sendMu.Lock()
	defer w.sendMu.Unlock()
	return len(b), w.conn.WriteMessage(websocket.TextMessage, b)
}

func (w *WSConn) Close() error {
	return w.conn.Close()
}

func (w *WSConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WSConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WSConn) SetDeadline(t time.Time) error {
	return nil
}

func (w *WSConn) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

func (w *WSConn) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}
