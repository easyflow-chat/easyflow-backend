package middleware

import (
	"easyflow-backend/src/common"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(module_name string, log_level common.LogLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.NewLogger(module_name, log_level)
		c.Set("logger", logger)
		c.Next()
	}
}
