package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type room struct {
	clients map[*client]bool
	join    chan *client
	leave   chan *client
	forward chan []byte
}

func newRoom() *room {
	return &room{
		clients: make(map[*client]bool),
		join:    make(chan *client),
		leave:   make(chan *client),
		forward: make(chan []byte),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			slog.Info("Add new client to room",
				"client.address", client.socket.RemoteAddr(), "clients", len(r.clients))

			r.clients[client] = true
		case client := <-r.leave:
			slog.Info("Remove client from room",
				"client.address", client.socket.RemoteAddr(), "clients", len(r.clients))

			delete(r.clients, client)
			client.close()
		case msg := <-r.forward:
			slog.Debug("Forward message to clients",
				"clients", len(r.clients))

			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		slog.Error("ServeHTTP", "error", err)
		return
	}

	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.read()
	client.read()
}
