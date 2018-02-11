package main

import (
	"net"
	"time"
)

type mockConn struct {
	Laddr  net.Addr
	Writes []string
}

func newMockConn() *mockConn {
	return &mockConn{Laddr: &net.TCPAddr{IP: net.IPv4(0, 0, 0, 1)}, Writes: make([]string, 0)}
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	m.Writes = append(m.Writes, string(b)) // keep track of all messages send to Write
	return len(b), err
}

func (m *mockConn) Close() error {
	panic("implement me")
}

func (m *mockConn) LocalAddr() net.Addr {
	panic("implement me")
}

func (m *mockConn) RemoteAddr() net.Addr {
	return m.Laddr
}

func (m *mockConn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
