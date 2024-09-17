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
		logger := common.NewLogger(os.Stdout, "AuthGuard", c)

		// Get access_token from cookies
		token, err := c.Cookie("access_token")
		if err != nil {
			logger.PrintfWarning("No access token found")
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
		logger := common.NewLogger(os.Stdout, "RefreshAuthGuard", c)
		// Get access_token from cookies
		token, err := c.Cookie("refresh_token")
		if err != nil {
			logger.PrintfWarning("No refresh token found")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		conf, ok := c.Get("config")
		if !ok {
			logger.PrintfError("Config not found in context")
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			})
			c.Abort()
			return
		}

		cfg, ok := conf.(*common.Config)
		if !ok {
			logger.PrintfError("Type assertion error")
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			})
			c.Abort()
			return
		}

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

		payload, err := ValidateToken(cfg, token)
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

		if err := db.(*gorm.DB).First(&database.UserKeys{}, "user_id = ? AND refresh_token = ?", payload.UserId, token).Error; err != nil {
			logger.PrintfWarning("Invalid refresh token")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidRefresh,
			})
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("refresh_token", "", -1, "/", cfg.BackendDomain, cfg.Stage != "development", false)

			c.Abort()
			return
		}

		c.Set("user", payload)
		c.Next()
	}
}
