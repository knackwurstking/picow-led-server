package ws

import "encoding/json"

type Request struct {
	Client  *Client `json:"-"`
	Command Command `json:"command"`
	Data    string  `json:"data"`
}

func NewRequest(c *Client, msg []byte) (*Request, error) {
	req := &Request{
		Client: c,
	}

	err := json.Unmarshal(msg, req)
	return req, err
}
