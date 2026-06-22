package apperror

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

func Conflict(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Err:     err,
	}
}

func ToHTTP(err error) (int, string) {
	var appError *AppError
	if errors.As(err, &appError) {
		return appError.Code, appError.Message
	}

	return http.StatusInternalServerError, "Internal Server Error"
}
