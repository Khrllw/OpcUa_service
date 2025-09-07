package cert_manager

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
)

// analyzeCertificate анализирует файл сертификата и выводит информацию о нем
func (cm *CertificateManager) analyzeCertificate(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		cm.logger.Error("Error reading file %s: %v", filePath, err)
		return fmt.Errorf("error reading file: %w", err)
	}

	cm.logger.Info("=== Analyzing certificate file: %s ===", filePath)
	cm.logger.Info("File size: %d bytes", len(data))

	// Hex dump of first 128 bytes
	cm.logger.Info("First 128 bytes (hex):\n%s", hex.Dump(data[:minim(128, len(data))]))

	// Try to parse as DER
	cert, err := x509.ParseCertificate(data)
	if err == nil {
		cm.logger.Info("File recognized as DER certificate")
		cm.printCertDetails(cert)
		return nil
	}

	// Try to parse as PEM
	block, _ := pem.Decode(data)
	if block != nil {
		cm.logger.Info("File contains PEM block: %s", block.Type)
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				cm.printCertDetails(cert)
				return nil
			}
		}
		cm.logger.Info("PEM block content (hex):\n%s", hex.Dump(block.Bytes[:minim(128, len(block.Bytes))]))
		return nil
	}

	cm.logger.Info("File is not recognized as DER or PEM certificate")
	cm.logger.Info("Full hex dump:\n%s", hex.Dump(data))
	return nil
}

// printCertDetails выводит основные поля сертификата
func (cm *CertificateManager) printCertDetails(cert *x509.Certificate) {
	cm.logger.Info("Certificate details:")
	cm.logger.Info("Subject: %s", cert.Subject)
	cm.logger.Info("Issuer: %s", cert.Issuer)
	cm.logger.Info("Serial Number: %s", cert.SerialNumber)
	cm.logger.Info("Valid From: %s", cert.NotBefore)
	cm.logger.Info("Valid Until: %s", cert.NotAfter)
	cm.logger.Info("Signature Algorithm: %s", cert.SignatureAlgorithm)
	cm.logger.Info("Public Key Algorithm: %s", cert.PublicKeyAlgorithm)

	if len(cert.URIs) > 0 {
		cm.logger.Info("URIs:")
		for _, uri := range cert.URIs {
			cm.logger.Info("- %s", uri)
		}
	}
}

// minim вспомогательная функция для minim
func minim(a, b int) int {
	if a < b {
		return a
	}
	return b
}
