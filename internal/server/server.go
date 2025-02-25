package server

import (
	"sync"

	"github.com/knackwurstking/picow-led-server/pkg/event"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
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

type Server struct {
	event *event.Event[*picow.Api]
	api   *picow.Api
	conns *connections
}

func NewServer(a *picow.Api, e *event.Event[*picow.Api]) *Server {
	return &Server{
		event: e,
		api:   a,
		conns: newConnections(),
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	// TODO: ...
}
