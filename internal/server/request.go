package server

import (
	"github.com/knackwurstking/picow-led-server/pkg/picow"
	"golang.org/x/net/websocket"
)

type Request struct {
	Conn    *websocket.Conn `json:"-"`
	Command Command         `json:"command"`
	// Data contains either an empty string or a json string from type `RequestData_*`
	Data string `json:"data"`
}

type (
	RequestData_PostApiDevice struct {
		picow.DeviceData
	}

	RequestData_PutApiDevice struct {
		picow.DeviceData
	}

	RequestData_DeleteApiDevice struct {
		Addr string `json:"addr"`
	}

	RequestData_PostApiDevicePins struct {
		Addr string     `json:"addr"`
		Pins picow.Pins `json:"pins"`
	}

	RequestData_PostApiDeviceColor struct {
		Addr  string      `json:"addr"`
		Color picow.Color `json:"color"`
	}
)
