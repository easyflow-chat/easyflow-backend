package user

import "time"

type CreateUserRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Name       string `json:"name" validate:"required,lte=50"`
	Password   string `json:"password" validate:"required,gte=8"`
	PublicKey  string `json:"publicKey" validate:"required"`
	PrivateKey string `json:"privateKey" validate:"required"`
	Iv         string `json:"iv" validate:"required"`
}

type CreateUserResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"update_at"`
	Email     string `json:"email"`
}

type GetUserResponse struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Bio        *string   `json:"bio"`
	Iv         string    `json:"iv"`
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private_key"`
}

type UpdateUserRequest struct {
	Name           *string `json:"name" validate:"omitempty,lte=50"`
	Bio            *string `json:"bio" validate:"omitempty,lte=1000"`
	ProfilePicture *string `json:"profile_picture" validate:"omitempty,lte=2048"`
}

type UpdateUserResponse struct {
	Id             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"update_at"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Bio            *string   `json:"bio"`
	ProfilePicture *string   `json:"profile_picture"`
}
