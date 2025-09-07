package errors

import "fmt"

type DBError struct {
	Message string
	Err     error
}

func NewDBError(message string, dbError error) *AppError {
	return &AppError{
		Code:         InternalServerErrorCode,
		Message:      message,
		Err:          dbError,
		IsUserFacing: false,
	}
}

func (e *DBError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *DBError) Unwrap() error {
	return e.Err
}
