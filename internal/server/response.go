package server

const (
	ResponseTypeDevices ResponseType = "devices"
	ResponseTypeDevice  ResponseType = "device"
	ResponseTypeError   ResponseType = "error"
)

type ResponseType string

type Response struct {
	Data any          `json:"data"`
	Type ResponseType `json:"type"`
}

func (r *Response) SetError(err error) {
	r.Type = ResponseTypeError
	r.Data = err.Error()
}

func (r *Response) SetData(t ResponseType, data any) {
	r.Type = t
	r.Data = data
}
