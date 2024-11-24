package common

import (
	"easyflow-backend/api"
	"fmt"

	"github.com/easyflow-chat/easyflow-backend/lib/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AnyStruct struct{}

func getPayload[T any](c *gin.Context) (*T, error) {
	var payload T

	if err := c.ShouldBind(&payload); err != nil && err.Error() != "EOF" {
		// Handle binding errors, but ignore io.EOF which occurs when the body is empty
		return nil, err
	}

	return &payload, nil
}

func getDatabse(c *gin.Context) (*gorm.DB, error) {
	raw_db, ok := c.Get("db")
	if !ok {
		return nil, fmt.Errorf("database not found in context")
	}

	db, ok := raw_db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("type assertion to *gorm.DB failed")
	}

	return db, nil
}

func getConfig(c *gin.Context) (*Config, error) {
	raw_cfg, ok := c.Get("config")
	if !ok {
		return nil, fmt.Errorf("Config not found in context")
	}

	cfg, ok := raw_cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("type assertion to *common.Config failed")
	}

	return cfg, nil
}

func getLogger(c *gin.Context) (*logger.Logger, error) {
	raw_logger, ok := c.Get("logger")
	if !ok {
		return nil, fmt.Errorf("Logger not found in context")
	}

	logger, ok := raw_logger.(*logger.Logger)
	if !ok {
		return nil, fmt.Errorf("type assertion to *logger.Logger failed")
	}

	return logger, nil
}

func SetupEndpoint[T any](c *gin.Context) (*T, *logger.Logger, *gorm.DB, *Config, []string) {
	var errors []error
	payload, err := getPayload[T](c)
	if err != nil {
		errors = append(errors, err)
	}

	db, err := getDatabse(c)
	if err != nil {
		errors = append(errors, err)
	}

	cfg, err := getConfig(c)
	if err != nil {
		errors = append(errors, err)
	}

	logger, err := getLogger(c)
	if err != nil {
		errors = append(errors, err)
	}

	var serializableErrors []string

	for _, e := range errors {
		if validationError, ok := e.(validator.ValidationErrors); ok {
			errArr := api.TranslateError(validationError)
			serializableErrors = append(serializableErrors, errArr...)
		} else {
			serializableErrors = append(serializableErrors, e.Error())
		}
	}

	return payload, logger, db, cfg, serializableErrors
}
