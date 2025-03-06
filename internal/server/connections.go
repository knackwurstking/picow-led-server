package server

import (
	"sync"

	"golang.org/x/net/websocket"
)

type Connections struct {
	conns map[*websocket.Conn]bool
	mutex *sync.Mutex
}

func NewConnections() *Connections {
	return &Connections{
		conns: map[*websocket.Conn]bool{},
		mutex: &sync.Mutex{},
	}
}

func (c *Connections) add(ws *websocket.Conn) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.conns[ws] = true
}

func (c *Connections) remove(ws *websocket.Conn) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	delete(c.conns, ws)
}

func (c *Connections) list() []*websocket.Conn {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	connections := make([]*websocket.Conn, 0)
	for conn := range c.conns {
		connections = append(connections, conn)
	}
	return connections
}
