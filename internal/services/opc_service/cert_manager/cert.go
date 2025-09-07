package cert_manager

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"
	"time"
)

type CertificateManager struct {
	logger *logging.Logger
}

func NewCertificateManager(logger *logging.Logger) interfaces.CertificateManagerService {
	return &CertificateManager{
		logger: logger.WithPrefix("CERT_MANAGER"),
	}
}

// DecodeCertificate принимает сертификат в виде байтов (DER или PEM)
func (cm *CertificateManager) DecodeCertificate(data []byte) (*x509.Certificate, error) {
	// Пробуем декодировать как PEM
	if block, _ := pem.Decode(data); block != nil && block.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("cannot parse PEM certificate: %w", err)
		}
		return cert, nil
	}

	// Пробуем как DER
	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("cannot parse DER certificate: %w", err)
	}

	return cert, nil
}

func (cm *CertificateManager) printEndpoints(resp []ua.EndpointDescription, selected *ua.EndpointDescription, policyID string) {
	cm.logger.Info("🔍 Available Endpoints:")

	for i, ep := range resp {
		cm.logger.Info("[%d] EndpointURL: %s", i, ep.EndpointURL)
		cm.logger.Info("    SecurityPolicy: %s", ep.SecurityPolicyURI)
		cm.logger.Info("    SecurityMode: %s", ep.SecurityMode)
		cm.logger.Info("    UserIdentityTokens:")
		for _, token := range ep.UserIdentityTokens {
			cm.logger.Info("      - Type: %s, PolicyID: %s, SecurityPolicy: %s",
				token.TokenType, token.PolicyID, token.SecurityPolicyURI)
		}
	}

	if selected != nil {
		cm.logger.Info("✅ Selected endpoint: %s", selected.EndpointURL)
		cm.logger.Info("   Security: %s + %s", selected.SecurityPolicyURI, selected.SecurityMode)
		cm.logger.Info("   Selected PolicyID: %s", policyID)
	} else {
		cm.logger.Info("✖ No endpoint supports certificate-based user authentication")
	}
}

func (cm *CertificateManager) SelectCertificateEndpoint(ctx context.Context, endpointURL string) (*ua.EndpointDescription, string, error) {
	resp, err := client.GetEndpoints(ctx, &ua.GetEndpointsRequest{
		EndpointURL: endpointURL,
	})
	if err != nil {
		return nil, "", err
	}

	var selected *ua.EndpointDescription
	var policyID string

	for _, ep := range resp.Endpoints {
		for _, token := range ep.UserIdentityTokens {
			if token.TokenType == ua.UserTokenTypeCertificate && selected == nil {
				selected = &ep
				policyID = token.PolicyID
			}
		}
	}

	//cm.printEndpoints(resp.Endpoints, selected, policyID)
	return selected, policyID, nil
}

func (cm *CertificateManager) PrintCertInfo(label string, cert *x509.Certificate) {
	cm.logger.Info("%s Cert: CN=%s, Valid: %s → %s\n",
		label,
		cert.Subject.CommonName,
		cert.NotBefore.Format(time.RFC3339),
		cert.NotAfter.Format(time.RFC3339),
	)
}

func (cm *CertificateManager) DecodeClientCredentials(certBytes, keyBytes []byte) ([]byte, *x509.Certificate, *rsa.PrivateKey) {
	// Парсим сертификат
	cert, err := cm.DecodeCertificate(certBytes)
	if err != nil {
		cm.logger.Error("Failed to load client certificate: %v", err)
		return nil, nil, nil
	}

	// Парсим приватный ключ
	key, err := cm.DecodePrivateKey(keyBytes)
	if err != nil {
		cm.logger.Error("Failed to load client private key: %v", err)
		return nil, nil, nil
	}

	return certBytes, cert, key
}

// TODO: Проверить потом без WithInsecureSkipVerify
func (cm *CertificateManager) BuildClientOptions(endpoint *ua.EndpointDescription, policyID string, certBytes []byte, key *rsa.PrivateKey) []client.Option {
	return []client.Option{
		client.WithClientCertificate(certBytes, key),
		client.WithX509Identity(certBytes, key),
		//client.WithTrustedCertificatesFile("certs/new_server/heopcua-rootca_khrllw--340595_2025aug06_132042.der"),
		client.WithSecurityPolicyURI(endpoint.SecurityPolicyURI, endpoint.SecurityMode),
		client.WithInsecureSkipVerify(),
	}
}
