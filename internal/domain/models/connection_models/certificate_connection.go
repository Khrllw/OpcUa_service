package models

import (
	"time"
)

// CertificateConnection конфигурация подключения OPC UA по сертификату
type CertificateConnection struct {
	EndpointURL  string
	Certificate  string
	Key          string
	Policy       string
	Mode         string
	Timeout      time.Duration
	Manufacturer string
	Model        string
}

func (*CertificateConnection) GetType() ConnectionTypeEnum {
	return ConnectionCertificate
}

func (cc *CertificateConnection) GetManufacturer() string {
	return cc.Manufacturer
}

func (cc *CertificateConnection) GetModel() string {
	return cc.Model
}

func (cc *CertificateConnection) GetEndpointURL() string {
	return cc.EndpointURL
}

func (ac *CertificateConnection) GetTimeout() time.Duration {
	return ac.Timeout
}
