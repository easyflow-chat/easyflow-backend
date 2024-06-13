package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTPayload struct {
	jwt.RegisteredClaims
	UserId      string     `json:"userId"`
	Email       string     `json:"email"`
	RefreshRand *uuid.UUID `json:"refreshRand"`
	IsPublic    bool       `json:"isPublic"`
}

type JWTPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
