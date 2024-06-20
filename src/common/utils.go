package common

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/enum"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StringPointer(s string) *string {
	return &s
}

// SetupEndpoint is a helper function to setup the endpoint with the necessary dependencies
// and return the payload, logger, db, and config
// PS: This function is not the preetiest but it gets the job done
func SetupEndpoint[T any](c *gin.Context, withPayload bool) (*T, *Logger, *gorm.DB, *Config, bool) {
	var payload T

	if withPayload {
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, api.ApiError{
				Code:    http.StatusBadRequest,
				Error:   enum.MalformedRequest,
				Details: err.Error(),
			})
			return nil, nil, nil, nil, false
		}

		if err := api.Validate.Struct(payload); err != nil {
			c.JSON(http.StatusBadRequest, api.ApiError{
				Code:    http.StatusBadRequest,
				Error:   enum.MalformedRequest,
				Details: api.TranslateError(err),
			})
			return nil, nil, nil, nil, false
		}
	}

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return nil, nil, nil, nil, false
	}

	cfg, ok := c.Get("config")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return nil, nil, nil, nil, false
	}

	raw_logger, ok := c.Get("logger")
	if !ok {
		panic("[User] Logger not found in context")
	}

	logger, ok := raw_logger.(*Logger)
	if !ok {
		panic("[User] Type assertion to *common.Logger failed")
	}

	return &payload, logger, db.(*gorm.DB), cfg.(*Config), true
}
