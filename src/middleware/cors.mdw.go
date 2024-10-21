package middleware

import (
	"easyflow-backend/src/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware(cfg *common.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentOrigin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Vary", "Origin")

		if currentOrigin == "" {
			c.Abort()
		}

		if currentOrigin != cfg.FrontendURL {
			c.Abort()
		}

		preflight := strings.ToUpper(c.Request.Method) == "OPTIONS"
		if preflight {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, Content-Length")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.FrontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if preflight {
			c.AbortWithStatus(200)
		}

		return
	}
}
