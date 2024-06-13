package user

import "github.com/golodash/galidator/v2"

type CreateUserRequest struct {
	Email      string `json:"email" binding:"required,max=255"`
	Name       string `json:"name" binding:"required,max=255"`
	Password   string `json:"password" binding:"required"`
	PublicKey  string `json:"public_key" binding:"required"`
	PrivateKey string `json:"private_key" binding:"required"`
	Iv         string `json:"iv" binding:"required,len=25" len:"$field must be $length characters"`
}

type CreateUserResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
	Email     string `json:"email"`
}

var (
	g = galidator.G()
)

var Validator = g.Validator(CreateUserRequest{}, galidator.Messages{
	"required": "$field is required",
	"max":      "$field must be less than $max characters",
	"len":      "$field must be $length characters",
})
