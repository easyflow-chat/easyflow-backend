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

var TOK_TTL = 60 * 60 * 24

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
			Code:  http.StatusUnauthorized,
			Error: enum.WrongCredentials,
		}
	}

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		logger.PrintfWarning("Wrong password for user with email: %s", payload.Email)
		return JWTPair{}, &api.ApiError{
			Code:  http.StatusUnauthorized,
			Error: enum.WrongCredentials,
		}
	}

	unique := uuid.New()
	expires := time.Now().Add(time.Duration(TOK_TTL) * time.Second)

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
			ExpiresAt: jwt.NewNumericDate(expires),
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
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	refreshToken, err := generateJwt(cfg, &refreshTokenPayload)

	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return JWTPair{}, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	//write refresh token to db
	entry := database.UserKeys{
		UserId:       user.Id,
		ExpiredAt:    expires,
		RefreshToken: refreshToken,
	}

	//check if user already has a refresh token
	if err := db.Where("user_id = ?", user.Id).First(&entry).Error; err != nil {
		logger.PrintfInfo("User with email: %s does not have a refresh token", payload.Email)
		if err := db.Create(&entry).Error; err != nil {
			logger.PrintfError("Error creating user key: %s", err)
			return JWTPair{}, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}
	} else {
		entry.RefreshToken = refreshToken
		entry.ExpiredAt = expires
		if err := db.Save(&entry).Error; err != nil {
			logger.PrintfError("Error updating user key: %s", err)
			return JWTPair{}, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}
	}

	return JWTPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func RefreshService(db *gorm.DB, cfg *common.Config, payload *JWTPayload, logger *common.Logger) (JWTPair, *api.ApiError) {
	//get user from db
	var user database.User
	if err := db.Where("id = ?", payload.UserId).First(&user).Error; err != nil {
		logger.PrintfInfo("User with id: %s not found", payload.UserId)
		return JWTPair{}, &api.ApiError{
			Code:  http.StatusUnauthorized,
			Error: enum.Unauthorized,
		}
	}

	unique := uuid.New()
	expires := time.Now().Add(time.Duration(TOK_TTL) * time.Second)

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

	/* _ := JWTPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Issuer:    "easyflow",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId:      user.Id,
		Email:       user.Email,
		RefreshRand: &unique,
		IsAccess:    false,
	} */

	accessToken, err := generateJwt(cfg, &accessTokenPayload)
	if err != nil {
		logger.PrintfError("Error generating jwt: %s", err)
		return JWTPair{}, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	return JWTPair{
		AccessToken:  accessToken,
		RefreshToken: "",
	}, nil
}
