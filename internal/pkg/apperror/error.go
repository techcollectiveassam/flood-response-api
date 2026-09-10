package apperror

import (
	"errors"
	"net/http"
)

const CodeInternalError = "internal_error"

type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	Code    string
	Status  int
	Message string
	Details []Detail
	cause   error
}

func New(status int, code, message string, cause ...error) *Error {
	e := &Error{Code: code, Status: status, Message: message}
	if len(cause) > 0 {
		e.cause = cause[0]
	}
	return e
}

func BadRequest(code, message string) *Error {
	return New(http.StatusBadRequest, code, message)
}

func NotFound(code, message string) *Error {
	return New(http.StatusNotFound, code, message)
}

func Conflict(code, message string) *Error {
	return New(http.StatusConflict, code, message)
}

func Internal(code, message string, cause ...error) *Error {
	return New(http.StatusInternalServerError, code, message, cause...)
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithDetails(details []Detail) *Error {
	e.Details = details
	return e
}

func (e *Error) Unwrap() error {
	return e.cause
}

func HTTPStatus(err error) int {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Status
	}
	return http.StatusInternalServerError
}

func Code(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeInternalError
}

func Message(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return "internal server error"
}

func Details(err error) []Detail {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Details
	}
	return nil
}
