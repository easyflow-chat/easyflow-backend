package api

import (
	"easyflow-backend/src/enum"
)

type ApiError struct {
	Code    int            `json:"code"`
	Error   enum.ErrorCode `json:"error"`
	Details interface{}    `json:"details,omitempty"`
}
