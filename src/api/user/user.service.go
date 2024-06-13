package user

import (
	"log"

	"easyflow-backend/src/api"
	"easyflow-backend/src/database"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, payload *CreateUserRequest) (*CreateUserResponse, *api.ApiError) {
	log.Println("Attempting to create user with email: ", payload.Email)
	var user database.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err == nil {
		return nil, &api.ErrUserAlreadyExists
	}

	//create a new user
	user = database.User{
		Email:      payload.Email,
		Name:       payload.Name,
		Password:   payload.Password,
		PublicKey:  payload.PublicKey,
		PrivateKey: payload.PrivateKey,
		Iv:         payload.Iv,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Println("Error creating user: ", err)
		return nil, &api.ErrFailedToCreateUser
	}

	return &CreateUserResponse{
		Id:        user.Id,
		CreatedAt: user.CreatedAt.String(),
		UpdateAt:  user.UpdatedAt.String(),
		Email:     user.Email,
	}, nil
}
