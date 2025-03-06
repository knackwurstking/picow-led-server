package server

import "github.com/knackwurstking/picow-led-server/pkg/picow"

type ResponseType string

const (
	ResponseTypeDevices ResponseType = "devices"
	ResponseTypeDevice  ResponseType = "device"
	ResponseTypeError   ResponseType = "error"
)

type Response[T string | []*picow.Device | *picow.Device] struct {
	Type ResponseType `json:"type"`
	Data T            `json:"data"`
}

type (
	ResponseError   = Response[string]
	ResponseDevice  = Response[*picow.Device]
	ResponseDevices = Response[[]*picow.Device]
)

func NewResponseError(e error) *ResponseError {
	return &ResponseError{
		Type: ResponseTypeError,
		Data: e.Error(),
	}
}

func NewResponseDevice(d *picow.Device) *ResponseDevice {
	return &ResponseDevice{
		Type: ResponseTypeDevice,
		Data: d,
	}
}

func NewResponseDevices(d []*picow.Device) *ResponseDevices {
	return &ResponseDevices{
		Type: ResponseTypeDevices,
		Data: d,
	}
}
