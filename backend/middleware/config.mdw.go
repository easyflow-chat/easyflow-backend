package middleware

import (
	"easyflow-backend/common"

	"github.com/gin-gonic/gin"
)

func ConfigMiddleware(cfg *common.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	}
}
