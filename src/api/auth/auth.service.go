package auth

import (
	"easyflow-backend/src/api"

	"gorm.io/gorm"
)

func LoginService(db *gorm.DB, payload *LoginRequest) (JWTPair, *api.ApiError) {

	return JWTPair{}, nil
}
