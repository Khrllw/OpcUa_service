package opc_connector

import (
	"encoding/base64"
	"fmt"
)

// Base64ToBytes преобразует base64-строку в []byte
func (oc *OpcConnector) Base64ToBytes(base64Str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return data, nil
}

// BytesToBase64 преобразует []byte в base64-строку
func (oc *OpcConnector) BytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
