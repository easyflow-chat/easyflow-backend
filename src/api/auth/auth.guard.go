package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access_token from cookies
		token, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:  http.StatusUnauthorized,
				Error: enum.Unauthorized,
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

		// Set user payload in context
		c.Set("user", payload)
		c.Next()
	}
}