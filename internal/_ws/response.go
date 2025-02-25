package ws

const (
	ResponseTypeDevices = "devices"
	ResponseTypeDevice  = "device"
	ResponseTypeError   = "error"
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

func (r *Response) Set(t ResponseType, data any) {
	r.Type = t
	r.Data = data
}
