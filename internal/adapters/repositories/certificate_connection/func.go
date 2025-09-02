package certificate_connection

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/pkg/errors"
)

func (r *CertificateConnectionRepositoryImpl) CreateCertificateConnection(cc entities.CertificateConnection) (uint, error) {
	op := "repo.CertificateConnection.CreateCertificateConnection"

	err := r.db.Clauses(clause.Returning{}).Create(&cc).Error
	if err != nil {
		return 0, errors.NewDBError(op, err)
	}

	return cc.ID, nil
}

func (r *CertificateConnectionRepositoryImpl) GetCertificateConnectionByID(id uint) (entities.CertificateConnection, error) {
	op := "repo.CertificateConnection.GetCertificateConnectionByID"

	var cc entities.CertificateConnection
	err := r.db.First(&cc, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.CertificateConnection{}, errors.NewDBError(op, errors.ErrNotFound)
		}
		return entities.CertificateConnection{}, errors.NewDBError(op, err)
	}

	return cc, nil
}

func (r *CertificateConnectionRepositoryImpl) UpdateCertificateConnection(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.CertificateConnection.UpdateCertificateConnection"

	var updated entities.CertificateConnection
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updated).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return updated.ID, nil
}

func (r *CertificateConnectionRepositoryImpl) DeleteCertificateConnection(id uint) error {
	op := "repo.CertificateConnection.DeleteCertificateConnection"

	result := r.db.Delete(&entities.CertificateConnection{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}
