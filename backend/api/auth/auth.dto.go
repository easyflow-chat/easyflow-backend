package auth

import (
	"github.com/easyflow-chat/easyflow-backend/lib/jwt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenResponse struct {
	jwt.JWTPair
	AccessTokenExpires int `json:"accessTokenExpires"`
}
