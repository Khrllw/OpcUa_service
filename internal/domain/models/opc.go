package models

import (
	"context"
	"github.com/awcullen/opcua/client"
	"sync"
	"time"
)

// ConnectionConfig конфигурация подключения OPC UA
type ConnectionConfig struct {
	EndpointURL string
	Certificate string
	Key         string
	Policy      string
	Mode        string
	Timeout     time.Duration
}

// ConnectorStats статистика менеджера подключений
type ConnectorStats struct {
	TotalConnections  int64
	FailedConnections int64
	ActiveConnections int64
	PoolSize          int64
}

// ConnectionInfo информация о подключении
type ConnectionInfo struct {
	Conn      *client.Client
	Ctx       context.Context
	Cancel    context.CancelFunc
	SessionID string
	Config    ConnectionConfig
	CreatedAt time.Time
	LastUsed  time.Time
	UseCount  int64
	IsHealthy bool
	Mu        sync.RWMutex
}
