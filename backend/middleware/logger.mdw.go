package middleware

import (
	"easyflow-backend/api"
	"easyflow-backend/common"
	"net/http"
	"os"

	"github.com/easyflow-chat/easyflow-backend/lib/logger"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(module_name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, ok := c.Get("config")
		if !ok {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    http.StatusInternalServerError,
				Error:   "ConfigError",
				Details: "Config not found in context",
			})
			c.Abort()
			return
		}

		config, ok := cfg.(*common.Config)
		if !ok {
			c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    http.StatusInternalServerError,
				Error:   "ConfigError",
				Details: "Config is not of type *common.Config",
			})
			c.Abort()
			return
		}

		c.Set("logger", logger.NewLogger(os.Stdout, module_name, logger.LogLevel(config.LogLevel), c.ClientIP()))
		c.Next()
	}
}
