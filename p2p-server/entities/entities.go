package entities
// just defining the different entities here, nothing much complicated

type (
	RegisterRequest struct {
		Name string `json:"name"`
	}

	RegisterResponse struct {
		Status bool   `json:"status"`
		Reason string `json:"reason"`
		Token  string `json:"token"`
	}

	UnregisterRequest struct {
		Token string `json:"token"`
	}

	UnregisterResponse struct {
		Status bool   `json:"status"`
		Reason string `json:"reason"`
	}

	ListRequest struct {
		Token string `json:"token"`
	}

	ListResponse struct {
		Status bool     `json:"status"`
		Reason string   `json:"reason"`
		Nodes  []string `json:"nodes"`
	}

	MessageRequest struct {
		Token string `json:"token"`
		Data  string `json:"data"`
	}

	MessageResponse struct {
		Status bool   `json:"status"`
		Reason string `json:"reason"`
	}

	SensorData struct {
		Name string            `json:"name"`
		Data map[string]string `json:"data"` // DataName : data
	}

	DeviceData struct {
		Name string                `json:"name"`
		Data map[string]SensorData `json:"data"` // SensorName : data
	}
// Gatewat data has an object of type device data which has an object of type sensor data
	GatewayData struct {
		Name string                `json:"name"`
		Data map[string]DeviceData `json:"data"` // DeviceName : data
	}
)
