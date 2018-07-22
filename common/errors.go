package common

import (
	"net/http"
)

// ApiError contains advanced information about an error, like a return status
// code to the desired action
type ApiError struct {
	Status  int
	Message string
}

// Error returns an error message
func (e *ApiError) Error() string {
	return e.Message
}

// Status returns the status code associated with this error
func (e *ApiError) Code() int {
	return e.Status
}

var (
	// ErrNotFound is an error associated with a 404 response
	ErrNotFound = &ApiError{http.StatusNotFound, "resource not found"}
	// ErrConflict is an error associated with a 409 response
	ErrConflict = &ApiError{http.StatusConflict, "duplicate resource id"}
	// ErrInternalServerError is an error associated with a 500 response
	ErrInternalServerError = &ApiError{http.StatusInternalServerError, "internal server error"}
	// ErrBadRequest returned when user submits an invalid request
	ErrBadRequest = &ApiError{http.StatusBadRequest, "bad request"}
)
