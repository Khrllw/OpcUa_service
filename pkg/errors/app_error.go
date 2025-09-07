package errors

import "fmt"

type AppError struct {
	Code         int    `json:"code"`    // HTTP статус код
	Message      string `json:"message"` // Сообщение для клиента
	Err          error  `json:"-"`       // Внутренняя ошибка, не для клиента
	IsUserFacing bool   `json:"-"`       // Флаг, указывающий, можно ли показывать `Err`
}

func NewAppError(httpCode int, message string, err error, isUserFacing bool) *AppError {
	return &AppError{
		Code:         httpCode,
		Message:      message,
		Err:          err,
		IsUserFacing: isUserFacing,
	}
}

func (ae *AppError) Unwrap() error {
	return ae.Err
}

func (ae *AppError) Error() string {
	if ae == nil {
		return ""
	}
	if ae.Err != nil {
		return fmt.Sprintf("%s (code: %d): %v", ae.Message, ae.Code, ae.Err)
	}
	return fmt.Sprintf("%s (code: %d)", ae.Message, ae.Code)
}
