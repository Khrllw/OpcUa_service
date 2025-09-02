package cert_manager

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"log"
	"opc_ua_service/internal/interfaces"
	"os"
	"time"
)

type CertificateManager struct {
}

func NewCertificateManager() interfaces.CertificateManagerService {
	return &CertificateManager{}
}

// Загрузка сертификата (DER или PEM)
func (cm *CertificateManager) LoadCertificate(certPath string) ([]byte, *x509.Certificate, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read cert file: %w", err)
	}

	if block, _ := pem.Decode(data); block != nil && block.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(block.Bytes)
		return block.Bytes, cert, err
	}

	cert, err := x509.ParseCertificate(data)
	return data, cert, err
}

// LoadCertificateBytes принимает сертификат в виде байтов (DER или PEM)
func (cm *CertificateManager) LoadCertificateBytes(data []byte) (*x509.Certificate, error) {
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

func (cm *CertificateManager) SelectCertificateEndpoint(ctx context.Context, endpointURL string) (*ua.EndpointDescription, string) {
	resp, err := client.GetEndpoints(ctx, &ua.GetEndpointsRequest{
		EndpointURL: endpointURL,
	})
	if err != nil {
		log.Fatalf("✖ GetEndpoints failed: %v", err)
	}

	fmt.Println("\n🔍 Available Endpoints:")
	var selected *ua.EndpointDescription
	var policyID string

	for i, ep := range resp.Endpoints {
		fmt.Printf("[%d] EndpointURL: %s\n", i, ep.EndpointURL)
		fmt.Printf("    SecurityPolicy: %s\n", ep.SecurityPolicyURI)
		fmt.Printf("    SecurityMode: %s\n", ep.SecurityMode)
		fmt.Printf("    UserIdentityTokens:\n")
		for _, token := range ep.UserIdentityTokens {
			fmt.Printf("      - Type: %s, PolicyID: %s, SecurityPolicy: %s\n",
				token.TokenType, token.PolicyID, token.SecurityPolicyURI)

			if token.TokenType == ua.UserTokenTypeCertificate && selected == nil {
				selected = &ep
				policyID = token.PolicyID
			}
		}
		fmt.Println()
	}

	if selected == nil {
		log.Fatal("✖ No endpoint supports certificate-based user authentication")
	}

	fmt.Printf("✅ Selected endpoint: %s\n", selected.EndpointURL)
	fmt.Printf("   Security: %s + %s\n", selected.SecurityPolicyURI, selected.SecurityMode)
	fmt.Printf("   Selected PolicyID: %s\n", policyID)

	return selected, policyID
}

func (cm *CertificateManager) PrintCertInfo(label string, cert *x509.Certificate) {
	fmt.Printf("✅ %s Cert: CN=%s, Valid: %s → %s\n",
		label,
		cert.Subject.CommonName,
		cert.NotBefore.Format(time.RFC3339),
		cert.NotAfter.Format(time.RFC3339),
	)
}

// --- вспомогательные функции ---

func (cm *CertificateManager) LoadClientCredentialsBase64(certBase64, keyBase64 string) ([]byte, *x509.Certificate, *rsa.PrivateKey) {
	// Декодируем сертификат из Base64
	certBytes, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		log.Fatalf("Failed to decode certificate Base64: %v", err)
	}

	// Декодируем ключ из Base64
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		log.Fatalf("Failed to decode key Base64: %v", err)
	}

	// Парсим сертификат
	cert, err := cm.LoadCertificateBytes(certBytes)
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	// Парсим приватный ключ
	key, err := cm.LoadPrivateKeyBytes(keyBytes)
	if err != nil {
		log.Fatalf("Failed to load client private key: %v", err)
	}

	return certBytes, cert, key
}

func (cm *CertificateManager) LoadClientCredentialsBytes(certBytes, keyBytes []byte) ([]byte, *x509.Certificate, *rsa.PrivateKey) {
	// Парсим сертификат
	cert, err := cm.LoadCertificateBytes(certBytes)
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	// Парсим приватный ключ
	key, err := cm.LoadPrivateKeyBytes(keyBytes)
	if err != nil {
		log.Fatalf("Failed to load client private key: %v", err)
	}

	return certBytes, cert, key
}

func (cm *CertificateManager) LoadClientCredentials(certPath, keyPath string) ([]byte, *x509.Certificate, *rsa.PrivateKey) {
	certBytes, cert, err := cm.LoadCertificate(certPath)
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	keyBytes, key, err := cm.LoadPrivateKey(keyPath)
	if err != nil {
		log.Fatalf("Failed to load client private key: %v", err)
	}

	// Кодируем в Base64 для безопасного отображения и копирования в JSON
	certBase64 := base64.StdEncoding.EncodeToString(certBytes)
	keyBase64 := base64.StdEncoding.EncodeToString(keyBytes)

	fmt.Println("____________________________________________________")
	fmt.Println("Load Client Cert (Base64):", certBase64)
	fmt.Println("____________________________________________________\n\n")
	fmt.Println("____________________________________________________")
	fmt.Println("Load Client Key (Base64):", keyBase64)
	fmt.Println("____________________________________________________")
	return certBytes, cert, key
}

func (cm *CertificateManager) LoadServerCertificate(certPath string) *x509.Certificate {
	_, cert, err := cm.LoadCertificate(certPath)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}
	return cert
}

// TODO: Проверить потом без WithInsecureSkipVerify
func (cm *CertificateManager) BuildClientOptions(endpoint *ua.EndpointDescription, policyID string, certBytes []byte, key *rsa.PrivateKey) []client.Option {
	return []client.Option{
		client.WithClientCertificate(certBytes, key),
		client.WithX509Identity(certBytes, key),
		//client.WithTrustedCertificatesFile("certs/new_server/heopcua-rootca_khrllw--340595_2025aug06_132042.der"),
		client.WithSecurityPolicyURI(endpoint.SecurityPolicyURI, endpoint.SecurityMode),
		client.WithInsecureSkipVerify(), // убрать в production
	}
}
