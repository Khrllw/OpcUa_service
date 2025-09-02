package cert_manager

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func (cm *CertificateManager) VerifyKeyMatchesCert(cert *x509.Certificate, key *rsa.PrivateKey) error {
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("certificate public key is not RSA")
	}

	if pub.N.Cmp(key.N) != 0 || pub.E != key.E {
		return fmt.Errorf("client certificate and private key do not match")
	}

	return nil
}

// LoadPrivateKeyBytes принимает приватный ключ в виде байтов (PEM/DER, PKCS#1/PKCS#8)
// Возвращает DER-байты и *rsa.PrivateKey
func (cm *CertificateManager) LoadPrivateKeyBytes(data []byte) (*rsa.PrivateKey, error) {
	// Пробуем как PEM
	if block, _ := pem.Decode(data); block != nil {
		switch block.Type {
		case "RSA PRIVATE KEY": // PKCS#1
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("cannot parse PKCS#1 private key: %w", err)
			}
			return key, nil

		case "PRIVATE KEY": // PKCS#8
			parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("cannot parse PKCS#8 private key: %w", err)
			}
			rsaKey, ok := parsedKey.(*rsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("not an RSA private key")
			}
			return rsaKey, nil

		default:
			return nil, fmt.Errorf("unsupported key type: %s", block.Type)
		}
	}

	// Пробуем как DER PKCS#8
	if parsedKey, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		rsaKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA private key")
		}
		return rsaKey, nil
	}

	// Пробуем как DER PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(data); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("cannot parse private key (unsupported format)")
}

// LoadPrivateKey Загрузка приватного ключа (PEM/DER, PKCS#1/PKCS#8)
func (cm *CertificateManager) LoadPrivateKey(keyPath string) ([]byte, *rsa.PrivateKey, error) {
	// Читаем файл
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read key file: %w", err)
	}

	// Пробуем декодировать как PEM
	if block, _ := pem.Decode(data); block != nil {
		switch block.Type {
		case "RSA PRIVATE KEY": // PKCS#1
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot parse PKCS#1 private key: %w", err)
			}
			return block.Bytes, key, nil

		case "PRIVATE KEY": // PKCS#8
			parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot parse PKCS#8 private key: %w", err)
			}
			rsaKey, ok := parsedKey.(*rsa.PrivateKey)
			if !ok {
				return nil, nil, fmt.Errorf("not an RSA private key")
			}
			return block.Bytes, rsaKey, nil

		default:
			return nil, nil, fmt.Errorf("unsupported key type: %s", block.Type)
		}
	}

	// Пробуем как DER PKCS#8
	if parsedKey, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		rsaKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, fmt.Errorf("not an RSA private key")
		}
		return data, rsaKey, nil
	}

	// Пробуем как DER PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(data); err == nil {
		return data, key, nil
	}

	return nil, nil, fmt.Errorf("cannot parse private key (unsupported format)")
}
