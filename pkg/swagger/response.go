package swagger

import (
	"opc_ua_service/internal/domain/models/connection"
)

type IDResponse struct {
	Status  string                `json:"status" example:"ok"`
	Message string                `json:"message" example:"Successfully connected"`
	Type    string                `json:"type" example:"object"`
	Data    connection.IDResponse `json:"data"`
}

type DisconnectResponse struct {
	Status  string                        `json:"status" example:"ok"`
	Message string                        `json:"message" example:"Successfully disconnected"`
	Type    string                        `json:"type" example:"object"`
	Data    connection.DisconnectResponse `json:"data"`
}

type CheckConnectionResponse struct {
	Status  string                            `json:"status" example:"ok"`
	Message string                            `json:"message" example:"Successfully get connection info"`
	Type    string                            `json:"type" example:"object"`
	Data    connection.ConnectionInfoResponse `json:"data"`
}

type GetConnectionPoolResponse struct {
	Status  string                            `json:"status" example:"ok"`
	Message string                            `json:"message" example:"Successfully get connection pool"`
	Type    string                            `json:"type" example:"object"`
	Data    connection.ConnectionPoolResponse `json:"data"`
}

type PollingResponse struct {
	Status string `json:"status" example:"ok"`

	Message string                     `json:"message" example:"Polling started/stopped for machine"`
	Type    string                     `json:"type" example:"object"`
	Data    connection.PollingResponse `json:"data"`
}
