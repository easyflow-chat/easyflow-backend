package auth

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type UserWithTokens struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Bio          *string   `json:"bio"`
	Iv           string    `json:"iv"`
	PublicKey    string    `json:"publicKey"`
	PrivateKey   string    `json:"privateKey"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}
