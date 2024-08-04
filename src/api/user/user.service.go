package user

import (
	"errors"
	"net/http"

	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/api/s3"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, payload *CreateUserRequest, cfg *common.Config, logger *common.Logger) (*CreateUserResponse, *api.ApiError) {
	logger.PrintfInfo("Attempting to create user with email: %s", payload.Email)
	var user database.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err == nil {
		logger.PrintfError("User with email: %s already exists", payload.Email)
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
	logger.PrintfInfo("Attempting to get user with email: %s", jwtPayload.Email)
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

func GetUserByEmail(db *gorm.DB, email string, logger *common.Logger) (bool, *api.ApiError) {
	logger.PrintfInfo("Attempting to find user with email: %s", email)
	var user database.User
	err := db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.PrintfInfo("No user with email: %s found", err)
		return false, nil
	}

	if err != nil {
		logger.PrintfInfo("An error occured while trying to find user: %s ", err)
		return false, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	return true, nil
}

func GenerateGetProfilePictureURL(db *gorm.DB, jwtPayload *auth.JWTPayload, logger *common.Logger, cfg *common.Config) (*string, *api.ApiError) {
	logger.PrintfInfo("Attempting to get profile picture for user with email: %s", jwtPayload.Email)
	var user database.User
	if err := db.Where("id = ?", jwtPayload.UserId).First(&user).Error; err != nil {
		logger.PrintfError("Error getting user: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.NotFound,
		}
	}

	imageURL, err := s3.GenerateDownloadURL(logger, cfg, cfg.ProfilePictureBucketName, user.Id, 60*60*24*7) // 1 week expiration time
	if err != nil {
		return nil, err
	}

	return imageURL, nil
}

func GenerateUploadProfilePictureURL(db *gorm.DB, jwtPayload *auth.JWTPayload, logger *common.Logger, cfg *common.Config) (*string, *api.ApiError) {
	logger.PrintfInfo("Attempting to create upload url for profile picture for user with email: %s", jwtPayload.Email)
	var user database.User
	if err := db.Where("id = ?", jwtPayload.UserId).First(&user).Error; err != nil {
		logger.PrintfError("Error getting user: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.NotFound,
		}
	}

	uploadURL, err := s3.GenerateUploadURL(logger, cfg, cfg.ProfilePictureBucketName, user.Id, 60*60)
	if err != nil {
		logger.PrintfError("Error uploading profile picture: %s", err.Error)
		return nil, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	return uploadURL, nil
}

func UpdateUser(db *gorm.DB, jwtPayload *auth.JWTPayload, payload *UpdateUserRequest, logger *common.Logger) (*UpdateUserResponse, *api.ApiError) {
	logger.PrintfInfo("Attempting to update user with email: %s", jwtPayload.UserId)
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

	if err := db.Update(user.Id, &user).Error; err != nil {
		logger.PrintfError("Error updating user: %s", err)
		return nil, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	return &UpdateUserResponse{
		Id:        user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Name:      user.Name,
		Bio:       user.Bio,
	}, nil
}

func DeleteUser(db *gorm.DB, jwtPayload *auth.JWTPayload, logger *common.Logger) *api.ApiError {
	logger.PrintfInfo("Attempting to delete user with email: %s", jwtPayload.UserId)
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
