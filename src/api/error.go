package api

import "net/http"

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrUserAlreadyExists  = ApiError{Code: http.StatusBadRequest, Message: "User already exists"}
	ErrFailedToCreateUser = ApiError{Code: http.StatusInternalServerError, Message: "Failed to create user"}
)
