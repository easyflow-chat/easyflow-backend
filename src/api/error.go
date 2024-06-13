package api

import (
	"net/http"
)

type ApiError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

var (
	//Internal errors
	ErrDatabaseConnection = ApiError{Code: http.StatusInternalServerError, Message: "Failed to connect to database"}
	//Users
	ErrUserAlreadyExists  = ApiError{Code: http.StatusBadRequest, Message: "User already exists"}
	ErrFailedToCreateUser = ApiError{Code: http.StatusInternalServerError, Message: "Failed to create user"}
)

// FieldError represents a detailed error message for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
