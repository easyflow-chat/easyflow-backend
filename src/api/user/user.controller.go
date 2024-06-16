package user

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("User"))
	r.POST("/signup", CreateUserController)
	r.GET("/", auth.AuthGuard(), GetUserController)
	r.GET("/profile-picture", auth.AuthGuard(), GetProfilePictureController)
	r.PUT("/", auth.AuthGuard(), UpdateUserController)
	r.DELETE("/", auth.AuthGuard(), DeleteUserController)
}

func CreateUserController(c *gin.Context) {
	var payload CreateUserRequest

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

	logger.Printf("Successfully validated request for creating user")

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

	user, err := CreateUser(db.(*gorm.DB), &payload, cfg.(*common.Config), logger)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	logger.Println("Responding with 200 status code for successful user creation")
	c.JSON(200, user)
}

func GetUserController(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	user, ok := val.(*auth.JWTPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	raw_logger, ok := c.Get("logger")
	if !ok {
		panic("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		panic("Type assertion to *common.Logger failed")
	}

	logger.Printf("Successfully validated request for getting user")

	db, ok := c.Get("db")
	if !ok {
		logger.PrintfError("Database not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	userFromDb, err := GetUserById(db.(*gorm.DB), user, logger)

	if err != nil {
		logger.PrintfError("Error getting user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}
	logger.Println("Responding with 200 status code for successful user retrieval")
	c.JSON(200, userFromDb)
}

func GetProfilePictureController(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		fmt.Println("User data not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	user, ok := val.(*auth.JWTPayload)
	if !ok {
		fmt.Println("Type assertion to JWTPayload failed")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	raw_logger, ok := c.Get("logger")
	if !ok {
		panic("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		panic("Type assertion to *common.Logger failed")
	}

	logger.Printf("Successfully validated request for getting profile picture")

	db, ok := c.Get("db")
	if !ok {
		logger.PrintfError("Database not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	pic, err := GetProfilePicture(db.(*gorm.DB), user, logger)

	if err != nil {
		logger.PrintfError("Error getting profile picture: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

	logger.Println("Responding with 200 status code for successful profile picture retrieval")
	c.JSON(200, pic)
}

func UpdateUserController(c *gin.Context) {
	val, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	user, ok := val.(*auth.JWTPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	var payload UpdateUserRequest
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
		panic("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		panic("Type assertion to *common.Logger failed")
	}

	logger.Printf("Successfully validated request for updating user")

	db, ok := c.Get("db")
	if !ok {
		logger.PrintfError("Database not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	updatedUser, err := UpdateUser(db.(*gorm.DB), user, &payload, logger)

	if err != nil {
		logger.PrintfError("Error updating user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

	logger.Println("Responding with 200 status code for successful user update")
	c.JSON(200, updatedUser)
}

func DeleteUserController(c *gin.Context) {
	val, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	user, ok := val.(*auth.JWTPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	raw_logger, ok := c.Get("logger")
	if !ok {
		panic("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		panic("Type assertion to *common.Logger failed")
	}

	logger.Printf("Successfully validated request for deleting user")

	db, ok := c.Get("db")
	if !ok {
		logger.PrintfError("Database not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	err := DeleteUser(db.(*gorm.DB), user, logger)

	if err != nil {
		logger.PrintfError("Error deleting user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

	logger.Println("Responding with 200 status code for successful user deletion")
	c.JSON(200, gin.H{
		"message": "User deleted",
	})
}
