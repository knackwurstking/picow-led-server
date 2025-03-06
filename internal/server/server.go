package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/knackwurstking/picow-led-server/pkg/event"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

const (
	EventNameChange = "change"
)

type Server struct {
	event *event.Event[*picow.Api]
	api   *picow.Api
	conns *Connections

	request chan *Request

	broadcastError   chan *ResponseError
	broadcastDevice  chan *ResponseDevice
	broadcastDevices chan *ResponseDevices

	mutexResponse   *sync.Mutex
	mutexApiDevices *sync.Mutex
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
		mutexResponse:    &sync.Mutex{},
		mutexApiDevices:  &sync.Mutex{},
	}
}

func (s *Server) StartResponseHandler() {
	for {
		select {
		case req := <-s.request:
			func() {
				switch req.Command {
				case CommandGetApiDevices:
					go Send(NewResponseDevices(s.api.Devices), s.mutexResponse, req.Conn)

				case CommandPostApiDevice:
					go func() {
						if req.Data == "" {
							return
						}

						reqData := RequestData_PostApiDevice{}
						if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
							Send(NewResponseError(err), s.mutexResponse, req.Conn)
							return
						}

						for _, d := range s.api.Devices {
							if d.Addr() == reqData.Server.Addr {
								Send(
									NewResponseError(
										fmt.Errorf(
											"device already exists, use \"%s\" command",
											CommandPutApiDevice,
										),
									),
									s.mutexResponse,
									req.Conn,
								)
								return
							}
						}

						s.api.Devices.Add(picow.NewDevice(reqData.DeviceData), s.mutexApiDevices)

						s.broadcastDevices <- NewResponseDevices(s.api.Devices)
						go s.event.Dispatch(EventNameChange, s.api)
					}()

				case CommandPutApiDevice:
					go func() {
						if req.Data == "" {
							return
						}

						reqData := RequestData_PutApiDevice{}
						if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
							Send(NewResponseError(err), s.mutexResponse, req.Conn)
							return
						}

						var device *picow.Device

						for _, d := range s.api.Devices {
							if d.Addr() == reqData.Server.Addr {
								device = d
								break
							}
						}

						if device == nil {
							Send(
								NewResponseError(
									fmt.Errorf(
										"device does not exist, use \"%s\" command",
										CommandPostApiDevice,
									),
								),
								s.mutexResponse,
								req.Conn,
							)
							return
						}

						device.SetDeviceData(reqData.DeviceData, s.mutexApiDevices)

						s.broadcastDevice <- NewResponseDevice(device)
						s.event.Dispatch(EventNameChange, s.api)
					}()

				case CommandDeleteApiDevice:
					go func() {
						if req.Data == "" {
							return
						}

						reqData := RequestData_DeleteApiDevice{}
						if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
							Send(NewResponseError(err), s.mutexResponse, req.Conn)
							return
						}

						var device *picow.Device
						for _, d := range s.api.Devices {
							if d.Addr() == reqData.Addr {
								device = d
								break
							}
						}
						if device == nil {
							return // No such device, just return without any error
						}

						if ok := s.api.Devices.Remove(device, s.mutexApiDevices); !ok {
							// Nothing to delete it seems
							return
						}

						s.broadcastDevices <- NewResponseDevices(s.api.Devices)
						s.event.Dispatch(EventNameChange, s.api)
					}()

				case CommandPostApiDevicePins:
					go func() {
						if req.Data == "" {
							return
						}

						reqData := &RequestData_PostApiDevicePins{}
						if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
							Send(NewResponseError(err), s.mutexResponse, req.Conn)
							return
						}

						// Checks
						if reqData.Pins == nil {
							Send(
								NewResponseError(fmt.Errorf("pins missing for %s", reqData.Addr)),
								s.mutexResponse,
								req.Conn,
							)
							return
						}

						var device *picow.Device

						for _, d := range s.api.Devices {
							if d.Addr() != reqData.Addr {
								continue
							}

							device = d

							if err := d.SetPins(reqData.Pins); err != nil {
								Send(NewResponseError(err), s.mutexResponse, req.Conn)
							}
						}

						// Handle response/broadcast
						if device == nil {
							Send(
								NewResponseError(
									fmt.Errorf("device %s not found", reqData.Addr),
								),
								s.mutexResponse,
								req.Conn,
							)
							return
						}

						s.broadcastDevice <- NewResponseDevice(device)
						s.event.Dispatch(EventNameChange, s.api)
					}()

				case CommandPostApiDeviceColor:
					go func() {
						if req.Data == "" {
							return
						}

						reqData := RequestData_PostApiDeviceColor{}
						if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
							Send(NewResponseError(err), s.mutexResponse, req.Conn)
							return
						}

						// Checks
						if reqData.Color == nil {
							Send(
								NewResponseError(fmt.Errorf("color missing for %s", reqData.Addr)),
								s.mutexResponse,
								req.Conn,
							)
							return
						}

						var device *picow.Device

						// Do stuff here
						for _, d := range s.api.Devices {
							if d.Addr() != reqData.Addr {
								continue
							}

							device = d

							if err := d.SetColor(reqData.Color); err != nil {
								Send(NewResponseError(err), s.mutexResponse, req.Conn)
							}
						}

						if device == nil {
							Send(
								NewResponseError(
									fmt.Errorf("device %s not found", reqData.Addr),
								),
								s.mutexResponse,
								req.Conn,
							)
							return
						}

						s.broadcastDevice <- NewResponseDevice(device)
					}()
				}
			}()

		case resp := <-s.broadcastError:
			go Send(resp, s.mutexResponse, s.conns.list()...)

		case resp := <-s.broadcastDevice:
			go Send(resp, s.mutexResponse, s.conns.list()...)

		case resp := <-s.broadcastDevices:
			go Send(resp, s.mutexResponse, s.conns.list()...)
		}
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	defer s.conns.remove(ws)
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

		r := &Request{Conn: ws}
		if json.Unmarshal(d, r) == nil {
			// NOTE: Ignore any error for now
			s.request <- r
		}
	}
}

func Send[T *ResponseError | *ResponseDevices | *ResponseDevice](data T, m *sync.Mutex, conns ...*websocket.Conn) {
	defer m.Unlock()
	m.Lock()

	slog.Debug("Sending data to clients", "data", data, "conns", conns)
	wg := &sync.WaitGroup{}

	if d, err := json.Marshal(data); err == nil {
		for _, c := range conns {
			wg.Add(1)

			go func() {
				defer wg.Done()

				slog.Debug("Writing data to a client", "data", string(d), "addr", c.RemoteAddr())
				c.SetWriteDeadline(time.Now().Add(time.Second * 5))
				if _, err := c.Write(d); err != nil {
					slog.Debug("Writing failed", "addr", c.RemoteAddr(), "error", err)
				}
			}()
		}

		wg.Wait()
	} else {
		slog.Error("Failed to marshal response", "error", err)
	}
}
