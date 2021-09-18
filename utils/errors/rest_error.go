package errors

import (
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewRestErr(message string, status int, error string) *RestErr {
	return &RestErr{message, status, error}
}

func NewBadRequestError(message string) *RestErr {
	return NewRestErr(message, http.StatusBadRequest, "bad_request")
}

func NewNotFoundError(message string) *RestErr {
	return NewRestErr(message, http.StatusNotFound, "not_found")
}

func NewInternalServerError(message string) *RestErr {
	return NewRestErr(message, http.StatusInternalServerError, "internal_server_error")
}
