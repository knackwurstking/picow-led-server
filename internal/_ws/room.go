package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

const (
	socketBufferSize = 1024
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	Api       *picow.Api
	Join      chan *Client
	Leave     chan *Client
	Handle    chan *Request
	Broadcast chan *Response

	// OnApiChange will trigger every time api data changed, excluded color changes
	OnApiChange func(a *picow.Api)

	clients map[*Client]bool

	mutexDevices *sync.Mutex
}

func NewRoom(api *picow.Api) *Room {
	return &Room{
		Api:          api,
		Join:         make(chan *Client),
		Leave:        make(chan *Client),
		Handle:       make(chan *Request),
		Broadcast:    make(chan *Response),
		OnApiChange:  nil,
		clients:      make(map[*Client]bool),
		mutexDevices: &sync.Mutex{},
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.Join:
			r.clients[client] = true

			slog.Debug(
				"Add a new client to the websocket room",
				"client.address", client.Socket.RemoteAddr(),
				"clients", len(r.clients),
			)

		case client := <-r.Leave:
			delete(r.clients, client)
			client.Close()

			slog.Debug(
				"Remove a client from the websocket room",
				"client.address", client.Socket.RemoteAddr(),
				"clients", len(r.clients),
			)

		case req := <-r.Handle:
			switch req.Command {
			case CommandGetApiDevices:
				go r.getApiDevices(req)
			case CommandPostApiDevice:
				go r.postApiDevice(req)
			case CommandPutApiDevice:
				go r.putApiDevice(req)
			case CommandDeleteApiDevice:
				go r.deleteApiDevice(req)
			case CommandPostApiDevicePins:
				go r.postApiDevicePins(req)
			case CommandPostApiDeviceColor:
				go r.postApiDeviceColor(req)
			}

		case resp := <-r.Broadcast:
			for c := range r.clients {
				go func(c *Client) {
					c.Response <- resp
				}(c)
			}
		}
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		slog.Error("ServeHTTP", "error", err)
		return
	}

	client := NewClient(socket, r)
	r.Join <- client
	defer func() {
		r.Leave <- client
	}()

	go client.Write()
	client.Read()
}

func (r *Room) getApiDevices(req *Request) {
	req.Client.Response <- &Response{
		Type: ResponseTypeDevices,
		Data: r.Api.Devices,
	}
}

func (r *Room) postApiDevice(req *Request) {
	if req.Data == "" {
		return
	}

	resp := &Response{}

	// Parse request data
	deviceData := picow.DeviceData{}
	if err := json.Unmarshal([]byte(req.Data), &deviceData); err != nil {
		resp.SetError(err)
		req.Client.Response <- resp
		return
	}

	// Checks
	for _, d := range r.Api.Devices {
		if d.Addr() == deviceData.Server.Addr {
			resp.SetError(
				fmt.Errorf(
					"device already exists, use \"%s\" command",
					CommandPutApiDevice,
				),
			)
			req.Client.Response <- resp
			return
		}
	}

	// Do stuff here
	r.Api.Devices.Add(picow.NewDevice(deviceData), r.mutexDevices)

	// Handle response/broadcast
	resp.Set(ResponseTypeDevices, r.Api.Devices)
	r.Broadcast <- resp

	if r.OnApiChange != nil {
		go r.OnApiChange(r.Api)
	}
}

func (r *Room) putApiDevice(req *Request) {
	if req.Data == "" {
		return
	}

	resp := &Response{}

	// Parse request data
	deviceData := picow.DeviceData{}
	if err := json.Unmarshal([]byte(req.Data), &deviceData); err != nil {
		resp.SetError(err)
		req.Client.Response <- resp
		return
	}

	// Checks
	var device *picow.Device

	for _, d := range r.Api.Devices {
		if d.Addr() == deviceData.Server.Addr {
			device = d
			break
		}
	}

	if device == nil {
		resp.SetError(
			fmt.Errorf(
				"device does not exist, use \"%s\" command",
				CommandPostApiDevice,
			),
		)
		req.Client.Response <- resp
		return
	}

	// Do stuff here
	device.SetDeviceData(deviceData, r.mutexDevices)

	// Handle response/broadcast
	resp.Set(ResponseTypeDevice, device)
	r.Broadcast <- resp

	if r.OnApiChange != nil {
		go r.OnApiChange(r.Api)
	}
}

func (r *Room) deleteApiDevice(req *Request) {
	if req.Data == "" {
		return
	}

	resp := &Response{}

	// Parse request data
	var reqData struct {
		Addr string `json:"addr"`
	}
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		resp.SetError(err)
		req.Client.Response <- resp
		return
	}

	// Checks
	var device *picow.Device
	for _, d := range r.Api.Devices {
		if d.Addr() == reqData.Addr {
			device = d
			break
		}
	}
	if device == nil {
		return // No such device, just return without any error
	}

	// Do stuff here
	if ok := r.Api.Devices.Remove(device, r.mutexDevices); !ok {
		// Nothing to delete it seems
		return
	}

	// Handle response/broadcast
	resp.Set(ResponseTypeDevices, r.Api.Devices)
	r.Broadcast <- resp

	if r.OnApiChange != nil {
		go r.OnApiChange(r.Api)
	}
}

func (r *Room) postApiDevicePins(req *Request) {
	if req.Data == "" {
		return
	}

	resp := &Response{}

	// Parse request data
	var reqData struct {
		Addr string     `json:"addr"`
		Pins picow.Pins `json:"pins"`
	}

	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		resp.SetError(err)
		req.Client.Response <- resp
		return
	}

	// Checks
	if reqData.Pins == nil {
		resp.SetError(
			fmt.Errorf("pins missing for %s", reqData.Addr),
		)
		req.Client.Response <- resp
		return
	}

	// Do stuff here
	for _, d := range r.Api.Devices {
		if d.Addr() != reqData.Addr {
			continue
		}

		if err := d.SetPins(reqData.Pins); err != nil {
			resp.SetError(err)
		} else {
			resp.Set(ResponseTypeDevice, d)
		}
	}

	// Handle response/broadcast
	if resp.Type == "" {
		resp.SetError(
			fmt.Errorf("device %s not found", reqData.Addr),
		)
		req.Client.Response <- resp
		return
	}
	r.Broadcast <- resp

	if r.OnApiChange != nil {
		go r.OnApiChange(r.Api)
	}
}

func (r *Room) postApiDeviceColor(req *Request) {
	if req.Data == "" {
		return
	}

	resp := &Response{}

	// Parse request data
	var reqData struct {
		Addr  string      `json:"addr"`
		Color picow.Color `json:"color"`
	}

	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		resp.SetError(err)
		req.Client.Response <- resp
		return
	}

	// Checks
	if reqData.Color == nil {
		resp.SetError(
			fmt.Errorf("color missing for %s", reqData.Addr),
		)
		req.Client.Response <- resp
		return
	}

	// Do stuff here
	for _, d := range r.Api.Devices {
		if d.Addr() != reqData.Addr {
			continue
		}

		if err := d.SetColor(reqData.Color); err != nil {
			resp.SetError(err)
		} else {
			resp.Set(ResponseTypeDevice, d)
		}
	}

	// Handle response/broadcast
	if resp.Type == "" {
		resp.SetError(
			fmt.Errorf("device %s not found", reqData.Addr),
		)
		req.Client.Response <- resp
		return
	}

	r.Broadcast <- resp
}
