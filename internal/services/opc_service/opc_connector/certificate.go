package opc_connector

import (
	"context"
	"github.com/awcullen/opcua/client"
	connectiion_models "opc_ua_service/internal/domain/models/connection_models"
	"time"
)

// ConnectWithCertificate Подключение по сертификату
func (oc *OpcConnector) ConnectWithCertificate(config connectiion_models.CertificateConnection) (*client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	certBytes, clientCert, clientKey := oc.certManager.LoadClientCredentialsBytes(config.Certificate, config.Key)
	if err := oc.certManager.VerifyKeyMatchesCert(clientCert, clientKey); err != nil {
		return nil, err
	}

	endpoint, policyID := oc.certManager.SelectCertificateEndpoint(ctx, config.EndpointURL)
	clientOpts := oc.certManager.BuildClientOptions(endpoint, policyID, certBytes, clientKey)

	conn, err := oc.createConnection(ctx, config.EndpointURL, clientOpts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// ConnectWithCertificateBytes Подключение по сертификату
func (oc *OpcConnector) ConnectWithCertificateBytes(config connectiion_models.CertificateConnection) (*client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	certBytes, clientCert, clientKey := oc.certManager.LoadClientCredentialsBytes([]byte(config.Certificate), []byte(config.Key))
	if err := oc.certManager.VerifyKeyMatchesCert(clientCert, clientKey); err != nil {
		return nil, err
	}

	endpoint, policyID := oc.certManager.SelectCertificateEndpoint(ctx, config.EndpointURL)
	clientOpts := oc.certManager.BuildClientOptions(endpoint, policyID, certBytes, clientKey)

	conn, err := oc.createConnection(ctx, config.EndpointURL, clientOpts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
