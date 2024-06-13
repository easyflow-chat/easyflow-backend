package auth

import "github.com/golang-jwt/jwt/v4"

var TOK_TTL = 60 * 60 * 24

type JWTPayload struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
	Email  string `json:"email"`
}

type JWTPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
