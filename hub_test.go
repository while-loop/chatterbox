package main

import (
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
)

func TestHub_Register(t *testing.T) {
	h := NewHub()
	go h.Start()
	c := newMockConn()
	h.Register(c)

	var wg sync.WaitGroup
	wg.Add(1)
	h.connChan <- func(conns map[string]net.Conn) {
		assert.Len(t, conns, 1)
		wg.Done()
	}

	wg.Wait()
	assert.Equal(t, "welcome!", c.Writes[0])
	assert.Len(t, c.Writes, 1)
}
