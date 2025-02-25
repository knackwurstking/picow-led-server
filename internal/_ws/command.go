package ws

const (
	CommandGetApiDevices      = Command("GET api.devices")
	CommandPostApiDevice      = Command("POST api.device")
	CommandPutApiDevice       = Command("PUT api.device")
	CommandDeleteApiDevice    = Command("DELETE api.device")
	CommandPostApiDevicePins  = Command("POST api.device.pins")
	CommandPostApiDeviceColor = Command("POST api.device.color")
)

type Command string
