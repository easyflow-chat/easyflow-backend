package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthEndpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("Auth"))
	r.POST("/login", LoginController)
	r.GET("/refresh", RefreshAuthGuard(), RefreshController)
}

func LoginController(c *gin.Context) {
	var payload LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    http.StatusBadRequest,
			Error:   enum.MalformedRequest,
			Details: err.Error(),
		})
		return
	}

	if err := api.Validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    http.StatusBadRequest,
			Error:   enum.MalformedRequest,
			Details: api.TranslateError(err),
		})
		return
	}

	raw_logger, ok := c.Get("logger")
	if !ok {
		log.Println("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		log.Println("Type assertion to *common.Logger failed")
	}

	logger.Printf("Successfully validated request for login")

	db, ok := c.Get("db")
	if !ok {
		logger.PrintfError("Database not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	cfg, ok := c.Get("config")
	if !ok {
		logger.PrintfError("Config not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	jwtPair, err := LoginService(db.(*gorm.DB), cfg.(*common.Config), &payload, logger)
	if err != nil {
		logger.PrintfError("Error logging in: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	c.SetCookie("access_token", jwtPair.AccessToken, TOK_TTL, "/", "", false, true)
	c.SetCookie("refresh_token", jwtPair.RefreshToken, TOK_TTL, "/", "", false, true)

	logger.Println("Responding with 200 status code for successful login")
	c.JSON(200, &gin.H{})
}

func RefreshController(c *gin.Context) {
	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	cfg, ok := c.Get("config")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	raw_logger, ok := c.Get("logger")
	if !ok {
		log.Println("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		log.Println("Type assertion to *common.Logger failed")
	}

	logger.Printf("Successfully validated request for token refresh")

	payload, ok := c.Get("user")
	if !ok {
		logger.PrintfError("User not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	jwtPair, err := RefreshService(db.(*gorm.DB), cfg.(*common.Config), payload.(*JWTPayload), logger)

	if err != nil {
		logger.PrintfError("Error refreshing token: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	c.SetCookie("access_token", jwtPair.AccessToken, TOK_TTL, "/", "", false, true)

	logger.Println("Responding with 200 status code for successful token refresh")
	c.JSON(200, &gin.H{})
}
