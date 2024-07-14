package web

import (
	"net/http"
)

type ErrorInterface interface {
	Code() string
	Description() string
	HttpStatusCode() int
}

type Error struct {
	code, description string
	httpStatusCode    int
}

func (err *Error) Code() string {
	return err.code
}

func (err *Error) Description() string {
	return err.description
}

func (err *Error) HttpStatusCode() int {
	return err.httpStatusCode
}

var ErrInternalServerError = func(description string) ErrorInterface {
	return &Error{"internal_server_error", description, http.StatusInternalServerError}
}

var ErrValidationFailed = func(desc string) ErrorInterface {
	return &Error{"bad_request", desc, http.StatusBadRequest}
}
