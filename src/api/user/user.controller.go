package user

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("User"))
	r.Use(middleware.RateLimiter(1, 4))
	r.POST("/signup", middleware.RateLimiter(1, 0), CreateUserController)
	r.GET("/", auth.AuthGuard(), GetUserController)
	r.GET("/exists/:email", UserExists)
	r.GET("/profile-picture", auth.AuthGuard(), GetProfilePictureController)
	r.POST("/profile-picture", auth.AuthGuard(), GenerateUploadProfilePictureURLController)
	r.PUT("/", auth.AuthGuard(), UpdateUserController)
	r.DELETE("/", auth.AuthGuard(), DeleteUserController)
}

func CreateUserController(c *gin.Context) {
	payload, logger, db, cfg, errors := common.SetupEndpoint[CreateUserRequest](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	logger.Printf("Successfully validated request for creating user")

	user, err := CreateUser(db, payload, cfg, logger)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	logger.Println("Responding with 200 status code for successful user creation")
	c.JSON(200, user)
}

func GetUserController(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

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

	logger.Printf("Successfully validated request for getting user")

	userFromDb, err := GetUserById(db, user, logger)

	if err != nil {
		logger.PrintfError("Error getting user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}
	logger.Println("Responding with 200 status code for successful user retrieval")
	c.JSON(200, userFromDb)
}

func GetProfilePictureController(c *gin.Context) {
	_, logger, db, cfg, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

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

	logger.Printf("Successfully validated request for getting profile picture")

	imageURL, err := GenerateGetProfilePicture(db, user, logger, cfg)

	if err != nil {
		logger.PrintfError("Error getting profile picture: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

	logger.Println("Responding with 200 status code for successful profile picture retrieval")
	c.JSON(200, imageURL)
}

func UserExists(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	email := c.Param("email")

	userInDb, err := GetUserByEmail(db, email, logger)

	if err != nil {
		logger.PrintfError("Error getting user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, userInDb)
}

func UpdateUserController(c *gin.Context) {
	payload, logger, db, _, errors := common.SetupEndpoint[UpdateUserRequest](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

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

	logger.Printf("Successfully validated request for updating user")

	updatedUser, err := UpdateUser(db, user, payload, logger)

	if err != nil {
		logger.PrintfError("Error updating user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

	logger.Println("Responding with 200 status code for successful user update")
	c.JSON(200, updatedUser)
}

func GenerateUploadProfilePictureURLController(c *gin.Context) {
	_, logger, db, cfg, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

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

	uploadURL, err := GenerateUploadProfilePictureURL(db, user, logger, cfg)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, uploadURL)
}

func DeleteUserController(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[CreateUserRequest](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}
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

	logger.Printf("Successfully validated request for deleting user")

	err := DeleteUser(db, user, logger)

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
