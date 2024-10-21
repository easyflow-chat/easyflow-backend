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
	IsAccess    bool       `json:"isAccess"`
}

type JWTPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
