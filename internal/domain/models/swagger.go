package models

type UUIDResponseSwagger struct {
	Status   string `json:"status" example:"ok"`
	Response struct {
		Message string       `json:"message" example:"Successfully connected"`
		Type    string       `json:"type" example:"object"`
		Data    UUIDResponse `json:"data"`
	} `json:"response"`
}

type DisconnectResponseSwagger struct {
	Status   string `json:"status" example:"ok"`
	Response struct {
		Message string             `json:"message" example:"Successfully disconnected"`
		Type    string             `json:"type" example:"object"`
		Data    DisconnectResponse `json:"data"`
	} `json:"response"`
}

type CheckConnectionResponseSwagger struct {
	Status   string `json:"status" example:"ok"`
	Response struct {
		Message string                 `json:"message" example:"Successfully get connection info"`
		Type    string                 `json:"type" example:"object"`
		Data    ConnectionInfoResponse `json:"data"`
	} `json:"response"`
}

type GetConnectionPoolResponseSwagger struct {
	Status   string `json:"status" example:"ok"`
	Response struct {
		Message string                 `json:"message" example:"Successfully get connection pool"`
		Type    string                 `json:"type" example:"object"`
		Data    ConnectionPoolResponse `json:"data"`
	} `json:"response"`
}

type PollingResponseSwagger struct {
	Status   string `json:"status" example:"ok"`
	Response struct {
		Message string          `json:"message" example:"Polling started/stopped for machine"`
		Type    string          `json:"type" example:"object"`
		Data    PollingResponse `json:"data"`
	} `json:"response"`
}
