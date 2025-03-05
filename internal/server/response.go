package server

import "github.com/knackwurstking/picow-led-server/pkg/picow"

type ResponseType string

const (
	ResponseTypeDevices ResponseType = "devices"
	ResponseTypeDevice  ResponseType = "device"
	ResponseTypeError   ResponseType = "error"
)

type Response[T string | []*picow.DeviceData | *picow.DeviceData] struct {
	Type ResponseType `json:"type"`
	Data T            `json:"data"`
}

type (
	ResponseError   = Response[string]
	ResponseDevice  = Response[*picow.DeviceData]
	ResponseDevices = Response[[]*picow.DeviceData]
)

func NewResponseError(d string) *ResponseError {
	return &ResponseError{
		Type: ResponseTypeError,
		Data: d,
	}
}

func NewResponseDevice(d *picow.DeviceData) *ResponseDevice {
	return &ResponseDevice{
		Type: ResponseTypeDevice,
		Data: d,
	}
}

func NewResponseDevices(d []*picow.DeviceData) *ResponseDevices {
	return &ResponseDevices{
		Type: ResponseTypeDevices,
		Data: d,
	}
}
