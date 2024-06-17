package user

import (
	"net/http"

	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, payload *CreateUserRequest, cfg *common.Config, logger *common.Logger) (*CreateUserResponse, *api.ApiError) {
	logger.Printf("Attempting to create user with email: %s", payload.Email)
	var user database.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err == nil {
		logger.Printf("User with email: %s already exists", payload.Email)
		return nil, &api.ApiError{
			Code:  http.StatusConflict,
			Error: enum.AlreadyExists,
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), cfg.SaltRounds)
	if err != nil {
		logger.PrintfError("Error hashing password: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	//create a new user
	user = database.User{
		Email:      payload.Email,
		Name:       payload.Name,
		Password:   string(password),
		PublicKey:  payload.PublicKey,
		PrivateKey: payload.PrivateKey,
		Iv:         payload.Iv,
	}

	if err := db.Create(&user).Error; err != nil {
		logger.PrintfError("Error creating user: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	logger.Printf("User created with email: %s", payload.Email)

	return &CreateUserResponse{
		Id:        user.Id,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	}, nil
}

func GetUserById(db *gorm.DB, jwtPayload *auth.JWTPayload, logger *common.Logger) (*GetUserResponse, *api.ApiError) {
	logger.Printf("Attempting to get user with email: %s", jwtPayload.Email)
	var user database.User
	if err := db.Where("id = ?", jwtPayload.UserId).First(&user).Error; err != nil {
		logger.PrintfError("Error getting user: %s", err)
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

func GetUserByEmail(db *gorm.DB, jwtPayload *auth.JWTPayload, logger *common.Logger) (*GetUserResponse, *api.ApiError) {
	logger.Printf("Attempting to get user with email: %s", jwtPayload.Email)
	var user database.User
	if err := db.Where("email = ?", jwtPayload.Email).First(&user).Error; err != nil {
		logger.PrintfError("Error getting user: %s", err)
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

func GetProfilePicture(db *gorm.DB, jwtPayload *auth.JWTPayload, logger *common.Logger) (*string, *api.ApiError) {
	logger.Printf("Attempting to get profile picture for user with email: %s", jwtPayload.Email)
	var user database.User
	if err := db.Where("id = ?", jwtPayload.UserId).First(&user).Error; err != nil {
		logger.PrintfError("Error getting user: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.NotFound,
		}
	}

	return user.ProfilePicture, nil
}

func UpdateUser(db *gorm.DB, jwtPayload *auth.JWTPayload, payload *UpdateUserRequest, logger *common.Logger) (*UpdateUserResponse, *api.ApiError) {
	logger.Printf("Attempting to update user with email: %s", jwtPayload.UserId)
	var user database.User
	if err := db.Where("id = ?", jwtPayload.UserId).First(&user).Error; err != nil {
		logger.PrintfError("Error getting user: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.NotFound,
		}
	}

	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.Bio != nil {
		user.Bio = payload.Bio
	}
	if payload.ProfilePicture != nil {
		user.ProfilePicture = payload.ProfilePicture
	}

	if err := db.Save(&user).Error; err != nil {
		logger.PrintfError("Error updating user: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	return &UpdateUserResponse{
		Id:             user.Id,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Email:          user.Email,
		Name:           user.Name,
		Bio:            user.Bio,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func DeleteUser(db *gorm.DB, jwtPayload *auth.JWTPayload, logger *common.Logger) *api.ApiError {
	logger.Printf("Attempting to delete user with email: %s", jwtPayload.UserId)
	var user database.User
	if err := db.Where("id = ?", jwtPayload.UserId).First(&user).Error; err != nil {
		logger.PrintfError("Error getting user: %s", err)
		return &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.NotFound,
		}
	}

	if err := db.Delete(&user).Error; err != nil {
		logger.PrintfError("Error deleting user: %s", err)
		return &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	return nil
}
