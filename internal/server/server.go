package server

import (
	"bytes"
	"io"
	"log/slog"

	"golang.org/x/net/websocket"

	"github.com/knackwurstking/picow-led-server/pkg/event"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
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
	defer func() {
		s.conns.remove(ws)
	}()
	defer ws.Close()

	// TODO: Add websocket to connections and start the main event loop, just a channel waiting for
	// event data
	s.conns.add(ws)

	var n int
	var err error
	var d []byte
	b := make([]byte, 1024)
outer:
	for {
		d = make([]byte, 0)

		for {
			// ws.SetReadDeadline(time.Now().Add(time.Second * 5))
			if n, err = ws.Read(b); err != nil || n == 0 {
				if err == io.EOF {
					slog.Debug("Got an end of file error")
					break outer
				}

				if err != nil {
					slog.Debug("Got an error", "error", err)
				}

				if n == 0 {
					slog.Debug("Got empty data from a client")
				}

				break
			} else {
				slog.Debug("Add buffer to data", "size", n, "char", string(b[:n]))
				d = append(d, b[:n]...)

				if b[n-1] == '\n' {
					slog.Debug("Detected a newline at the end")
					break
				}
			}
		}

		d = bytes.TrimRight(d, "\n")

		// TODO: Handle request, do somethings with data here
		slog.Debug("Got some data from a client", "size", len(d), "data", string(d))
	}
}
