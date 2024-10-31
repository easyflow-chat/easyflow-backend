package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, logger, _, cfg, errs := common.SetupEndpoint[any](c)
		if errs != nil {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    http.StatusInternalServerError,
				Error:   enum.ApiError,
				Details: errs,
			})
			c.Abort()
			return
		}

		// Get access_token from cookies
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			logger.PrintfWarning("Error while getting access token cookie: %s", err.Error())
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if accessToken == "" {
			logger.PrintfWarning("No access token provided")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		// Validate token
		payload, err := ValidateToken(cfg, accessToken)
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

		// Set user payload in context
		c.Set("user", payload)
		c.Next()
	}
}

func RefreshAuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, logger, db, cfg, errs := common.SetupEndpoint[any](c)
		if errs != nil {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    http.StatusInternalServerError,
				Error:   enum.ApiError,
				Details: errs,
			})
			c.Abort()
			return
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			logger.PrintfWarning("Error while getting refresh token cookie: %s", err.Error())
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if refreshToken == "" {
			logger.PrintfWarning("No refresh token provided")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		token, err := ValidateToken(cfg, refreshToken)
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

		if token.Issuer != "easyflow" {
			logger.PrintfWarning("Invalid issuer")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
			})
			c.Abort()
			return
		}

		if err := db.First(&database.UserKeys{}, "user_id = ? AND random = ?", token.UserId, token.RefreshRand).Error; err != nil {
			logger.PrintfWarning("Invalid refresh token")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidRefresh,
			})

			c.Abort()
			return
		}

		c.Set("user", token)
		c.Next()
	}
}
