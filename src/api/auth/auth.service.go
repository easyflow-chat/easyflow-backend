package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"
	"fmt"
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

func LoginService(db *gorm.DB, cfg *common.Config, payload *LoginRequest) (JWTPair, *api.ApiError) {
	var user database.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		fmt.Println("Error finding user: ", err)
		return JWTPair{}, &api.ApiError{
			Code:  http.StatusUnauthorized,
			Error: enum.WrongCredentials,
		}
	}

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		fmt.Println("Error comparing password: ", err)
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
		IsPublic:    true,
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
		IsPublic:    false,
	}

	accessToken, err := generateJwt(cfg, &accessTokenPayload)

	if err != nil {
		fmt.Println("Error generating jwt: ", err)
		return JWTPair{}, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	refreshToken, err := generateJwt(cfg, &refreshTokenPayload)

	if err != nil {
		fmt.Println("Error generating jwt: ", err)
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
		if err := db.Create(&entry).Error; err != nil {
			fmt.Println("Error creating user key: ", err)
			return JWTPair{}, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}
	} else {
		entry.RefreshToken = refreshToken
		entry.ExpiredAt = expires
		if err := db.Save(&entry).Error; err != nil {
			fmt.Println("Error updating user key: ", err)
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
