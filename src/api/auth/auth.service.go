package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func generateJwt(cfg *common.Config, payload *JWTPayload) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := claims.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(cfg *common.Config, token string) (*JWTPayload, error) {
	claims := &JWTPayload{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func LoginService(db *gorm.DB, cfg *common.Config, payload *LoginRequest, logger *common.Logger) (JWTPair, *api.ApiError) {
	var user database.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		logger.PrintfWarning("User with email: %s not found", payload.Email)
		return JWTPair{}, &api.ApiError{
			Code:    http.StatusUnauthorized,
			Error:   enum.WrongCredentials,
			Details: err,
		}
	}

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		logger.PrintfWarning("Wrong password for user with email: %s", payload.Email)
		return JWTPair{}, &api.ApiError{
			Code:    http.StatusUnauthorized,
			Error:   enum.WrongCredentials,
			Details: err,
		}
	}

	unique := uuid.New()
	expires := time.Now().Add(time.Duration(cfg.JwtExpirationTime) * time.Second)
	refreshExpires := time.Now().Add(time.Duration(cfg.RefreshExpirationTime) * time.Second)

	accessTokenPayload := JWTPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "easyflow",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId:      user.Id,
		Email:       user.Email,
		RefreshRand: &unique,
		IsAccess:    true,
	}

	refreshTokenPayload := JWTPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpires),
			Issuer:    "easyflow",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId:      user.Id,
		Email:       user.Email,
		RefreshRand: &unique,
		IsAccess:    false,
	}

	accessToken, err := generateJwt(cfg, &accessTokenPayload)

	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return JWTPair{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	refreshToken, err := generateJwt(cfg, &refreshTokenPayload)

	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return JWTPair{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	//write refresh token to db
	entry := database.UserKeys{
		Id:           unique.String(),
		ExpiredAt:    refreshExpires,
		RefreshToken: refreshToken,
		UserId:       user.Id,
	}

	if err := db.Create(&entry).Error; err != nil {
		logger.PrintfError("Error updating user key: %s", err)
		return JWTPair{}, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	return JWTPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func RefreshService(db *gorm.DB, cfg *common.Config, payload *JWTPayload, logger *common.Logger) (*JWTPair, *api.ApiError) {
	//get user from db
	var user database.User
	if err := db.First(&user, "id = ?", payload.UserId).Error; err != nil {
		logger.PrintfError("Could not get user with id: %s", payload.UserId)
		return nil, &api.ApiError{
			Code:    http.StatusUnauthorized,
			Error:   enum.Unauthorized,
			Details: err,
		}
	}

	unique := uuid.New()
	expires := time.Now().Add(time.Duration(cfg.JwtExpirationTime) * time.Second)

	accessTokenPayload := JWTPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "easyflow",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId:      user.Id,
		Email:       user.Email,
		RefreshRand: &unique,
		IsAccess:    true,
	}

	accessToken, err := generateJwt(cfg, &accessTokenPayload)
	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return nil, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	return &JWTPair{
		AccessToken:  accessToken,
		RefreshToken: "",
	}, nil
}

func LogoutService(db *gorm.DB, payload *JWTPayload, logger *common.Logger) *api.ApiError {
	if err := db.Delete(&database.UserKeys{}, payload.RefreshRand).Error; err != nil {
		logger.PrintfError("Could not delete Refresh Token with id: %s", payload.RefreshRand)
		return &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	return nil
}
