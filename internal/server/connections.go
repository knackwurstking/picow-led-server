package server

import (
	"sync"

	"golang.org/x/net/websocket"
)

type connections struct {
	conns map[*websocket.Conn]bool
	mutex *sync.Mutex
}

func newConnections() *connections {
	return &connections{
		conns: map[*websocket.Conn]bool{},
		mutex: &sync.Mutex{},
	}
}

func (c connections) add(ws *websocket.Conn) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.conns[ws] = true
}

func (c connections) remove(ws *websocket.Conn) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	delete(c.conns, ws)
}
