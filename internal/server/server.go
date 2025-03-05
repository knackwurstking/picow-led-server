package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/knackwurstking/picow-led-server/pkg/event"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

type Server struct {
	event            *event.Event[*picow.Api]
	api              *picow.Api
	conns            *Connections
	request          chan *Request
	broadcastError   chan *ResponseError
	broadcastDevice  chan *ResponseDevice
	broadcastDevices chan *ResponseDevices
	responseMutex    *sync.Mutex
}

func NewServer(a *picow.Api, e *event.Event[*picow.Api]) *Server {
	return &Server{
		event:            e,
		api:              a,
		conns:            NewConnections(),
		request:          make(chan *Request),
		broadcastError:   make(chan *ResponseError),
		broadcastDevice:  make(chan *ResponseDevice),
		broadcastDevices: make(chan *ResponseDevices),
		responseMutex:    &sync.Mutex{},
	}
}

func (s *Server) StartResponseHandler() {
	for {
		select {
		case req := <-s.request:
			func() {
				// TODO: Handle requests here...
				switch req.Command {
				case CommandGetApiDevices:
					respond(NewResponseDevices(s.api.Devices), req.Conn)

				case CommandPostApiDevice:
					// ...

				case CommandPutApiDevice:
					// ...

				case CommandDeleteApiDevice:
					// ...

				case CommandPostApiDevicePins:
					// ...

				case CommandPostApiDeviceColor:
					// ...
				}
			}()

		case resp := <-s.broadcastError:
			func() {
				defer s.responseMutex.Unlock()
				s.responseMutex.Lock()
				respond(resp, s.conns.list()...)
			}()

		case resp := <-s.broadcastDevice:
			func() {
				defer s.responseMutex.Unlock()
				s.responseMutex.Lock()
				respond(resp, s.conns.list()...)
			}()

		case resp := <-s.broadcastDevices:
			func() {
				defer s.responseMutex.Unlock()
				s.responseMutex.Lock()
				respond(resp, s.conns.list()...)
			}()
		}
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	defer func() {
		s.conns.remove(ws)
	}()
	defer ws.Close()

	s.conns.add(ws)
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	var (
		n   int
		err error
		d   []byte
	)

	b := make([]byte, 1024)

main:
	for {
		d = make([]byte, 0)

	buffer:
		for {
			// ws.SetReadDeadline(time.Now().Add(time.Second * 5))
			if n, err = ws.Read(b); err != nil || n == 0 {
				if err == io.EOF {
					slog.Debug("Got an end of file error")
					break main
				}

				if err != nil {
					slog.Debug("Got an error", "error", err)
				}

				if n == 0 {
					slog.Debug("Got empty data from a client")
				}

				break buffer
			} else {
				slog.Debug("Add buffer to data", "size", n, "char", string(b[:n]))
				d = append(d, b[:n]...)

				if b[n-1] == '\n' {
					slog.Debug("Detected a newline at the end")
					break buffer
				}
			}
		}

		d = bytes.TrimRight(d, "\n")
		slog.Debug("Got some data from a client", "size", len(d), "data", string(d))

		r := &Request{}
		if json.Unmarshal(d, r) == nil {
			// NOTE: Ignore any error for now
			s.request <- r
		}
	}
}

func respond[T *ResponseError | *ResponseDevices | *ResponseDevice](data T, conns ...*websocket.Conn) {
	if d, err := json.Marshal(data); err == nil {
		for _, c := range conns {
			go func() {
				if _, err := c.Write(d); err != nil {
					slog.Debug("Writing failed", "addr", c.RemoteAddr(), "error", err)
				}
			}()
		}
	} else {
		slog.Error("Failed to marshal response", "error", err)
	}
}
