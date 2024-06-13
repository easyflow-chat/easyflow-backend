package user

type CreateUserRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Name       string `json:"name" validate:"required,lte=50"`
	Password   string `json:"password" validate:"required,gte=8"`
	PublicKey  string `json:"public_key" validate:"required"`
	PrivateKey string `json:"private_key" validate:"required"`
	Iv         string `json:"iv" validate:"required,len=25"`
}

type CreateUserResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
	Email     string `json:"email"`
}
