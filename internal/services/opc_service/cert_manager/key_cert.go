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

// LoadPrivateKey Загрузка приватного ключа (PEM/DER, PKCS#1/PKCS#8)
func (cm *CertificateManager) LoadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read key file: %w", err)
	}

	if block, _ := pem.Decode(data); block != nil {
		switch block.Type {
		case "RSA PRIVATE KEY":
			return x509.ParsePKCS1PrivateKey(block.Bytes)
		case "PRIVATE KEY":
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			rsaKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("not an RSA private key")
			}
			return rsaKey, nil
		default:
			return nil, fmt.Errorf("unsupported key type: %s", block.Type)
		}
	}

	if key, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
	}
	if key, err := x509.ParsePKCS1PrivateKey(data); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("cannot parse private key (unsupported format)")
}
