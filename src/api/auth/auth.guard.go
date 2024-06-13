package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

		// Validate token
		payload, err := ValidateToken(cfg.(*common.Config), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if payload.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.ExpiredToken,
			})
			c.Abort()
			return
		}

		if payload.Issuer != "easyflow" {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if !payload.IsAccess {
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

		db, ok := c.Get("db")
		if !ok {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			})
			c.Abort()
			return
		}

		payload, err := ValidateToken(cfg.(*common.Config), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if payload.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.ExpiredToken,
			})
			c.Abort()
			return
		}

		if payload.Issuer != "easyflow" {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if payload.IsAccess {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		if err := db.(*gorm.DB).Where("refresh_token = ?", token).First(&database.UserKeys{}).Error; err != nil {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidRefresh,
			})
			c.Abort()
			return
		}

		c.Set("user", payload)
		c.Next()
	}
}
