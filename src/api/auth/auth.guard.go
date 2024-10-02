package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"
	"errors"
	"net/http"
	"strings"

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

		// Get access_token from header
		var token string
		header := c.GetHeader("Authorization")
		if header != "" {
			token = strings.Split(header, "Bearer ")[1]
			if token == "" {
				logger.PrintfWarning("No access token found")
				c.JSON(http.StatusUnauthorized, api.ApiError{
					Code:  http.StatusUnauthorized,
					Error: enum.InvalidCookie,
				})
				c.Abort()
				return
			}
		} else {
			logger.PrintfWarning("No access token found")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		// Validate token
		payload, err := ValidateToken(cfg, token)
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
		payload, logger, db, cfg, errs := common.SetupEndpoint[RefreshTokenRequest](c)
		if errs != nil {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    http.StatusInternalServerError,
				Error:   enum.ApiError,
				Details: errs,
			})
			c.Abort()
			return
		}

		token, err := ValidateToken(cfg, payload.RefreshToken)
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

		if token.IsAccess {
			logger.PrintfWarning("Invalid token type")
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.InvalidCookie,
			})
			c.Abort()
			return
		}

		if err := db.First(&database.UserKeys{}, "user_id = ? AND refresh_token = ?", token.UserId, token).Error; err != nil {
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
