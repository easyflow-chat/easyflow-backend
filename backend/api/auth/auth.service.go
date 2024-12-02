package auth

import (
	"easyflow-backend/api"
	"easyflow-backend/api/utils"
	"easyflow-backend/common"
	"easyflow-backend/enum"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/easyflow-chat/easyflow-backend/lib/database"
	"github.com/easyflow-chat/easyflow-backend/lib/jwt"
	"github.com/easyflow-chat/easyflow-backend/lib/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginService(db *gorm.DB, cfg *common.Config, payload *LoginRequest, logger *logger.Logger) (jwt.JWTPair, database.User, *api.ApiError) {
	var user database.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		logger.PrintfWarning("User with email: %s not found", payload.Email)
		return jwt.JWTPair{}, database.User{}, &api.ApiError{
			Code:    http.StatusUnauthorized,
			Error:   enum.WrongCredentials,
			Details: err,
		}
	}

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		logger.PrintfWarning("Wrong password for user with email: %s", payload.Email)
		return jwt.JWTPair{}, database.User{}, &api.ApiError{
			Code:    http.StatusUnauthorized,
			Error:   enum.WrongCredentials,
			Details: err,
		}
	}

	random := uuid.New()
	expires := time.Now().Add(time.Duration(cfg.JwtExpirationTime) * time.Second)
	refreshExpires := time.Now().Add(time.Duration(cfg.RefreshExpirationTime) * time.Second)

	accessTokenPayload := jwt.CreateTokenPayload(user.ID, random.String(), expires, false)

	refreshTokenPayload := jwt.CreateTokenPayload(user.ID, random.String(), refreshExpires, true)

	accessToken, err := jwt.GenerateJwt[jwt.JWTTokenPayload](cfg.JwtSecret, accessTokenPayload)
	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return jwt.JWTPair{}, database.User{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	refreshToken, err := jwt.GenerateJwt[jwt.JWTTokenPayload](cfg.JwtSecret, refreshTokenPayload)
	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return jwt.JWTPair{}, database.User{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	//write refresh token to db
	entry := database.UserKeys{
		Random:    random.String(),
		ExpiredAt: refreshExpires,
		UserID:    user.ID,
	}

	if err := db.Create(&entry).Error; err != nil {
		logger.PrintfError("Error updating user key: %s", err)
		return jwt.JWTPair{}, user, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	if user.ProfilePicture == nil {
		utils.GenerateNewProfilePictureUrl(logger, cfg, db, &user)
	} else {
		expired := false

		pictureUrl, err := url.Parse(*user.ProfilePicture)
		if err != nil {
			expired = true
		}

		query := pictureUrl.Query()
		issuedAt, err := time.Parse(time.RFC3339, query.Get("X-Amz-Date"))
		if err != nil {
			expired = true
		}
		expiryTime, err := strconv.ParseInt(query.Get("X-Amz-Expires"), 10, 64)
		if err != nil {
			expired = true
		}

		if issuedAt.Add(time.Duration(expiryTime) * time.Second).After(time.Now()) {
			expired = true
		}

		if expired {
			utils.GenerateNewProfilePictureUrl(logger, cfg, db, &user)
		}

	}

	logger.Printf("Logged in user: %s", user.ID)

	return jwt.JWTPair{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, user, nil
}

func RefreshService(db *gorm.DB, cfg *common.Config, payload *jwt.JWTTokenPayload, logger *logger.Logger) (jwt.JWTPair, *api.ApiError) {
	//get user from db
	var user database.User
	if err := db.First(&user, "id = ?", payload.UserID).Error; err != nil {
		logger.PrintfWarning("Could not get user with id: %s", payload.UserID)
		return jwt.JWTPair{}, &api.ApiError{
			Code:    http.StatusUnauthorized,
			Error:   enum.Unauthorized,
			Details: err,
		}
	}

	random := uuid.New().String()
	expires := time.Now().Add(time.Duration(cfg.JwtExpirationTime) * time.Second)
	refreshExpires := time.Now().Add(time.Duration(cfg.RefreshExpirationTime) * time.Second)

	accessTokenPayload := jwt.CreateTokenPayload(user.ID, random, expires, false)

	refreshTokenPayload := jwt.CreateTokenPayload(user.ID, random, refreshExpires, true)

	accessToken, err := jwt.GenerateJwt(cfg.JwtSecret, &accessTokenPayload)
	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return jwt.JWTPair{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	refreshToken, err := jwt.GenerateJwt(cfg.JwtSecret, &refreshTokenPayload)
	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return jwt.JWTPair{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	//write refresh token random to db
	if err := db.Model(database.UserKeys{}).Where(
		&database.UserKeys{
			UserID: payload.UserID,
			Random: payload.RefreshRand,
		},
	).Updates(
		database.UserKeys{
			Random:    random,
			ExpiredAt: refreshExpires,
		}).Error; err != nil {
		logger.PrintfError("Error updating user key with user id: %s and random: %s due to: %s", payload.UserID, payload.RefreshRand, err)
		return jwt.JWTPair{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	logger.Printf("Refreshed token for user with id: %s and random: %s. New random: %s", payload.UserID, payload.RefreshRand, random)

	return jwt.JWTPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func LogoutService(db *gorm.DB, payload *jwt.JWTTokenPayload, logger *logger.Logger) *api.ApiError {
	if err := db.Delete(&database.UserKeys{UserID: payload.UserID, Random: payload.RefreshRand}).Error; err != nil {
		logger.PrintfError("Could not delete Refresh Token with random: %s and user id: %s", payload.RefreshRand, payload.UserID)
		return &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	logger.Printf("Successfully ended session for user with id: %s and random: %s", payload.UserID, payload.RefreshRand)

	return nil
}
