package user

import (
	"fmt"
	"log"
	"net/http"

	"easyflow-backend/src/api"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, payload *CreateUserRequest) (*CreateUserResponse, *api.ApiError) {
	log.Println("Attempting to create user with email: ", payload.Email)
	var user database.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err == nil {
		return nil, &api.ApiError{
			Code:  http.StatusConflict,
			Error: enum.AlreadyExists,
		}
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
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	return &CreateUserResponse{
		Id:        user.Id,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	}, nil
}

func GetUserById(db *gorm.DB, id *string) (*GetUserResponse, *api.ApiError) {
	fmt.Println("Attempting to get user with email: ", id)
	var user database.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	return &GetUserResponse{
		Id:         user.Id,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Email:      user.Email,
		Name:       user.Name,
		Bio:        user.Bio,
		Iv:         user.Iv,
		PublicKey:  user.PublicKey,
		PrivateKey: user.PrivateKey,
	}, nil
}
