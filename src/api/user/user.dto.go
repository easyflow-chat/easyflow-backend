package user

type CreateUserRequest struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Iv         string `json:"iv"`
}

/*
{
	"email": "test@test.com",
	"name": "test",
	"password": "test",
	"public_key": "test",
	"private_key": "test",
	"iv": "test"
}
*/

type CreateUserResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
	Email     string `json:"email"`
}
