package errors

import (
	"errors"
	"fmt"
)

const (
	InternalServerError = "internal server error"
	BadRequest          = "bad request"
	NotFound            = "not_found"

	BadRequestErrorCode     = 400
	InvalidDataCode         = 400
	ForbiddenErrorCode      = 403
	InternalServerErrorCode = 500
	NotFoundErrorCode       = 404
)

var (
	ErrEmptyAction  = errors.New("action did not affect the data")
	ErrDataNotFound = errors.New("data not found")
	ErrEmptyData    = errors.New("empty data")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal error")
)

func Is(err any, err2 error) bool {
	if e, ok := err.(error); ok {
		return errors.Is(e, err2)
	}
	return false
}

var ErrNotFound = errors.New("not found")

func NewNotFoundError(message string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, message)
}

// NewInternalError создает ошибку внутреннего сервера
func NewInternalError(op string, message string, err error) *AppError {
	return &AppError{
		Code:         InternalServerErrorCode,
		Message:      fmt.Sprintf("%s: %s", op, message),
		Err:          fmt.Errorf("%w: %v", ErrInternal, err),
		IsUserFacing: false,
	}
}

// NewForbiddenError создает ошибку доступа (может пригодиться в будущем)
func NewForbiddenError(op string, message string) *AppError {
	return &AppError{
		Code:         ForbiddenErrorCode,
		Message:      fmt.Sprintf("%s: %s", op, message),
		Err:          ErrForbidden,
		IsUserFacing: true,
	}
}
