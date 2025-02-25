package ws

import (
	"encoding/json"
	"log/slog"

	"github.com/gorilla/websocket"
)

type Client struct {
	Socket   *websocket.Conn
	Response chan *Response
	Room     *Room
}

func NewClient(s *websocket.Conn, r *Room) *Client {
	return &Client{
		Socket:   s,
		Response: make(chan *Response),
		Room:     r,
	}
}

func (c *Client) Read() {
	defer c.Socket.Close()

	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			slog.Debug(
				"Error while reading a message from a client",
				"client.address", c.Socket.RemoteAddr(),
				"error", err,
			)
			return
		}

		req, err := NewRequest(c, msg)
		if err != nil {
			slog.Warn("Parsing request failed", "error", err)
			continue
		}

		slog.Debug(
			"Got a message from a client",
			"client.address", c.Socket.RemoteAddr(),
			"request.command", req.Command,
			"request.data.length", len(req.Data),
		)

		c.Room.Handle <- req
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()

	for resp := range c.Response {
		data, err := json.Marshal(resp)
		if err != nil {
			slog.Warn("Marshal response failed",
				"client.address", c.Socket.RemoteAddr(),
				"error", err)

			resp.SetError(err)
			data, _ = json.Marshal(resp)
		}

		if err := c.Socket.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}

func (c *Client) Close() {
	defer func() {
		if r := recover(); r != nil {
			slog.Debug(
				"Recovered while closing a client",
				"client.address", c.Socket.RemoteAddr(),
				"error", r,
			)
		}
	}()

	c.Socket.Close()
	close(c.Response)
}
