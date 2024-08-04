package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access_token from cookies
		token, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		// Get config from context
		cfg, ok := c.Get("config")
		if !ok {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			})
			c.Abort()
			return
		}

		logger := common.NewLogger(os.Stdout, "AuthGuard")

		// Validate token
		payload, err := ValidateToken(cfg.(*common.Config), token)
		if err != nil {
			logger.PrintfError("Error validating token: %s", err.Error())
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, api.ApiError{
					Code:    http.StatusUnauthorized,
					Error:   enum.ExpiredToken,
					Details: err,
				})
			} else {
				c.JSON(http.StatusInternalServerError, api.ApiError{
					Code:    http.StatusInternalServerError,
					Error:   enum.ApiError,
					Details: err,
				})
			}
			c.Abort()
			return
		}

		if payload.Issuer != "easyflow" {
			logger.PrintfWarning("Invalid issuer")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if !payload.IsAccess {
			logger.PrintfWarning("Invalid token type")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		// Set user payload in context
		c.Set("user", payload)
		c.Next()
	}
}

func RefreshAuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access_token from cookies
		token, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		cfg, ok := c.Get("config")
		if !ok {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			})
			c.Abort()
			return
		}

		logger := common.NewLogger(os.Stdout, "RefreshAuthGuard")

		db, ok := c.Get("db")
		if !ok {
			logger.PrintfError("Database not found in context")
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			})
			c.Abort()
			return
		}

		payload, err := ValidateToken(cfg.(*common.Config), token)
		if err != nil {
			logger.PrintfError("Error validating token: %s", err.Error())
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, api.ApiError{
					Code:  http.StatusUnauthorized,
					Error: enum.ExpiredRefreshToken,
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if payload.Issuer != "easyflow" {
			logger.PrintfWarning("Invalid issuer")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if payload.IsAccess {
			logger.PrintfWarning("Invalid token type")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		var key database.UserKeys

		if err := db.(*gorm.DB).Where("Id = ?", payload.UserId).First(&key).Error; err != nil || key.RefreshToken != token {
			logger.PrintfWarning("Invalid refresh token")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidRefresh,
			})
			c.SetCookie("refresh_token", "", -1, "/", "", cfg.(common.Config).Stage != "development", true)

			c.Abort()
			return
		}

		c.Set("user", payload)
		c.Next()
	}
}
