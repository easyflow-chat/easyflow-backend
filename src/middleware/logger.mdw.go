package middleware

import (
	"easyflow-backend/src/common"
	"os"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(module_name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", common.NewLogger(os.Stdout, module_name))
		c.Next()
	}
}
