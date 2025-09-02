package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"opc_ua_service/internal/domain/models"
)

// ClientAPI определяет интерфейс для взаимодействия с сервисом
type ClientAPI interface {
	CreateConnection(ctx context.Context, req *models.ConnectionRequest) (*models.UUIDResponseSwagger, *http.Response, error)
	GetConnectionPool(ctx context.Context) (*models.GetConnectionPoolResponseSwagger, *http.Response, error)
	DeleteConnection(ctx context.Context, req *models.UUIDRequest) (*models.DisconnectResponseSwagger, *http.Response, error)
	CheckConnection(ctx context.Context, req *models.CheckConnectionRequest) (*models.CheckConnectionResponseSwagger, *http.Response, error)
	StartPolling(ctx context.Context, req *models.UUIDRequest) (*models.PollingResponseSwagger, *http.Response, error)
	StopPolling(ctx context.Context, req *models.UUIDRequest) (*models.PollingResponseSwagger, *http.Response, error)
}

// Client реализует интерфейс ClientAPI
type Client struct {
	service *ClientService
}

// NewClient создает нового клиента
func NewClient(host string) ClientAPI {
	return &Client{
		service: NewClientService(host),
	}
}

// CreateConnection создает новое подключение
func (c *Client) CreateConnection(ctx context.Context, req *models.ConnectionRequest) (*models.UUIDResponseSwagger, *http.Response, error) {
	const endpoint = "/api/v1/connect"

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := c.service.createRequestJSONWithContext(ctx, http.MethodPost, endpoint, nil, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	body, httpResp, err := c.service.doRequest(httpReq)
	if err != nil {
		return nil, httpResp, err
	}

	var resp models.UUIDResponseSwagger
	if err := json.Unmarshal(body, &resp); err != nil {

		return nil, httpResp, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	fmt.Println("Raw response body:", string(body))
	return &resp, httpResp, nil
}

// GetConnectionPool возвращает пул активных соединений
func (c *Client) GetConnectionPool(ctx context.Context) (*models.GetConnectionPoolResponseSwagger, *http.Response, error) {
	const endpoint = "/api/v1/connect"

	httpReq, err := c.service.createRequestJSONWithContext(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	body, httpResp, err := c.service.doRequest(httpReq)
	if err != nil {
		return nil, httpResp, err
	}

	var resp models.GetConnectionPoolResponseSwagger
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, httpResp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, httpResp, nil
}

// DeleteConnection отключает сессию по UUID
func (c *Client) DeleteConnection(ctx context.Context, req *models.UUIDRequest) (*models.DisconnectResponseSwagger, *http.Response, error) {
	const endpoint = "/api/v1/connect"

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := c.service.createRequestJSONWithContext(ctx, http.MethodDelete, endpoint, nil, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	body, httpResp, err := c.service.doRequest(httpReq)
	if err != nil {
		return nil, httpResp, err
	}

	var resp models.DisconnectResponseSwagger
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, httpResp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, httpResp, nil
}

// CheckConnection проверяет состояние соединения
func (c *Client) CheckConnection(ctx context.Context, req *models.CheckConnectionRequest) (*models.CheckConnectionResponseSwagger, *http.Response, error) {
	const endpoint = "api/v1/connect/check"

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := c.service.createRequestJSONWithContext(ctx, http.MethodPost, endpoint, nil, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	body, httpResp, err := c.service.doRequest(httpReq)
	if err != nil {
		return nil, httpResp, err
	}

	var resp models.CheckConnectionResponseSwagger
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, httpResp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, httpResp, nil
}

// StartPolling запускает опрос OPC UA по UUID станка
func (c *Client) StartPolling(ctx context.Context, req *models.UUIDRequest) (*models.PollingResponseSwagger, *http.Response, error) {
	const endpoint = "/api/v1/polling/start"

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := c.service.createRequestJSONWithContext(ctx, http.MethodGet, endpoint, nil, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	body, httpResp, err := c.service.doRequest(httpReq)
	if err != nil {
		return nil, httpResp, err
	}

	var resp models.PollingResponseSwagger
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, httpResp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, httpResp, nil
}

// StopPolling останавливает опрос OPC UA по UUID станка
func (c *Client) StopPolling(ctx context.Context, req *models.UUIDRequest) (*models.PollingResponseSwagger, *http.Response, error) {
	const endpoint = "/api/v1/polling/stop"

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := c.service.createRequestJSONWithContext(ctx, http.MethodGet, endpoint, nil, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	body, httpResp, err := c.service.doRequest(httpReq)
	if err != nil {
		return nil, httpResp, err
	}

	var resp models.PollingResponseSwagger
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, httpResp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, httpResp, nil
}
