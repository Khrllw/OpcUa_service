package anonymous_connection

import "gorm.io/gorm"

type AnonymousConnectionRepositoryImpl struct {
	db *gorm.DB
}

func NewAnonymousConnectionRepository(db *gorm.DB) *AnonymousConnectionRepositoryImpl {
	return &AnonymousConnectionRepositoryImpl{db: db}
}
