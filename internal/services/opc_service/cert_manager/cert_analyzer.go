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
		return fmt.Errorf("error reading file: %w", err)
	}

	fmt.Printf("\n=== Analyzing certificate file: %s ===\n", filePath)
	fmt.Printf("File size: %d bytes\n", len(data))

	// Hex dump of first 128 bytes
	fmt.Println("\nFirst 128 bytes (hex):")
	fmt.Println(hex.Dump(data[:min(128, len(data))]))

	// Try to parse as DER
	cert, err := x509.ParseCertificate(data)
	if err == nil {
		fmt.Println("\nFile recognized as DER certificate")
		cm.printCertDetails(cert)
		return nil
	}

	// Try to parse as PEM
	block, _ := pem.Decode(data)
	if block != nil {
		fmt.Printf("\nFile contains PEM block: %s\n", block.Type)
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				cm.printCertDetails(cert)
				return nil
			}
		}
		fmt.Println("PEM block content (hex):")
		fmt.Println(hex.Dump(block.Bytes[:min(128, len(block.Bytes))]))
		return nil
	}

	fmt.Println("\nFile is not recognized as DER or PEM certificate")
	fmt.Println("Full hex dump:")
	fmt.Println(hex.Dump(data))
	return nil
}

// printCertDetails выводит основные поля сертификата
func (cm *CertificateManager) printCertDetails(cert *x509.Certificate) {
	fmt.Println("\nCertificate details:")
	fmt.Printf("Subject: %s\n", cert.Subject)
	fmt.Printf("Issuer: %s\n", cert.Issuer)
	fmt.Printf("Serial Number: %s\n", cert.SerialNumber)
	fmt.Printf("Valid From: %s\n", cert.NotBefore)
	fmt.Printf("Valid Until: %s\n", cert.NotAfter)
	fmt.Printf("Signature Algorithm: %s\n", cert.SignatureAlgorithm)
	fmt.Printf("Public Key Algorithm: %s\n", cert.PublicKeyAlgorithm)

	if len(cert.URIs) > 0 {
		fmt.Println("\nURIs:")
		for _, uri := range cert.URIs {
			fmt.Printf("- %s\n", uri)
		}
	}
}
