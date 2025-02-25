package server

import (
	"github.com/knackwurstking/picow-led-server/pkg/event"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
	"golang.org/x/net/websocket"
)

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
