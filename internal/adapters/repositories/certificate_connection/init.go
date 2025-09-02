package certificate_connection

import "gorm.io/gorm"

type CertificateConnectionRepositoryImpl struct {
	db *gorm.DB
}

func NewCertificateConnectionRepository(db *gorm.DB) *CertificateConnectionRepositoryImpl {
	return &CertificateConnectionRepositoryImpl{db: db}
}
