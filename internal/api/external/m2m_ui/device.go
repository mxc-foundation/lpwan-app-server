package m2m_ui

type DeviceServerAPI struct {
	serviceName string
}

func NewDeviceServerAPI() *DeviceServerAPI {
	return &DeviceServerAPI{serviceName: "device"}
}

