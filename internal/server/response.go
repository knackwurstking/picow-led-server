package server

import "github.com/knackwurstking/picow-led-server/pkg/picow"

const (
	ResponseTypeDevices ResponseType = "devices"
	ResponseTypeDevice  ResponseType = "device"
	ResponseTypeError   ResponseType = "error"
)

type ResponseType string

type (
	ResponseError   = Response[string]
	ResponseDevice  = Response[*picow.Device]
	ResponseDevices = Response[[]*picow.DeviceData]
)

type Response[T string | []*picow.DeviceData | *picow.Device] struct {
	Type ResponseType `json:"type"`
	Data any          `json:"data"`
}

func (r *Response[T]) SetError(err error) {
	r.Type = ResponseTypeError
	r.Data = err.Error()
}

func (r *Response[T]) SetData(t ResponseType, data any) {
	r.Type = t
	r.Data = data
}
