package middleware

import (
	"easyflow-backend/src/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware(cfg *common.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentOrigin := c.Request.Header.Get("Origin")

		if currentOrigin == "" {
			return
		}

		if currentOrigin != cfg.FrontendURL {
			return
		}

		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.FrontendURL)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

			c.Next()
		}
	}
}
