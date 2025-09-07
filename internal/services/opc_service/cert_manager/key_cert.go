package cert_manager

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// VerifyKeyMatchesCert проверяет, что публичный ключ сертификата соответствует приватному ключу
func (cm *CertificateManager) VerifyKeyMatchesCert(cert *x509.Certificate, key *rsa.PrivateKey) error {
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		cm.logger.Error("Certificate public key is not RSA")
		return fmt.Errorf("certificate public key is not RSA")
	}

	if pub.N.Cmp(key.N) != 0 || pub.E != key.E {
		cm.logger.Error("Client certificate and private key do not match")
		return fmt.Errorf("client certificate and private key do not match")
	}

	cm.logger.Info("Certificate and private key match successfully")
	return nil
}

// DecodePrivateKey принимает приватный ключ в виде байтов (PEM/DER, PKCS#1/PKCS#8)
// Возвращает *rsa.PrivateKey
func (cm *CertificateManager) DecodePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// Пробуем как PEM
	if block, _ := pem.Decode(data); block != nil {
		cm.logger.Info("PEM block detected: %s", block.Type)
		switch block.Type {
		case "RSA PRIVATE KEY": // PKCS#1
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				cm.logger.Error("Cannot parse PKCS#1 private key: %v", err)
				return nil, fmt.Errorf("cannot parse PKCS#1 private key: %w", err)
			}
			cm.logger.Info("PKCS#1 private key parsed successfully")
			return key, nil

		case "PRIVATE KEY": // PKCS#8
			parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				cm.logger.Error("Cannot parse PKCS#8 private key: %v", err)
				return nil, fmt.Errorf("cannot parse PKCS#8 private key: %w", err)
			}
			rsaKey, ok := parsedKey.(*rsa.PrivateKey)
			if !ok {
				cm.logger.Error("Parsed PKCS#8 key is not an RSA private key")
				return nil, fmt.Errorf("not an RSA private key")
			}
			cm.logger.Info("PKCS#8 private key parsed successfully")
			return rsaKey, nil

		default:
			cm.logger.Error("Unsupported PEM key type: %s", block.Type)
			return nil, fmt.Errorf("unsupported key type: %s", block.Type)
		}
	}

	// Пробуем как DER PKCS#8
	if parsedKey, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		rsaKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			cm.logger.Error("Parsed DER PKCS#8 key is not an RSA private key")
			return nil, fmt.Errorf("not an RSA private key")
		}
		cm.logger.Info("DER PKCS#8 private key parsed successfully")
		return rsaKey, nil
	}

	// Пробуем как DER PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(data); err == nil {
		cm.logger.Info("DER PKCS#1 private key parsed successfully")
		return key, nil
	}

	cm.logger.Error("Cannot parse private key (unsupported format)")
	return nil, fmt.Errorf("cannot parse private key (unsupported format)")
}
