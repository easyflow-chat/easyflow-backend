package user

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
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
	r.GET("/upload-profile-picture", auth.AuthGuard(), GenerateUploadProfilePictureURLController)
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

	user, err := CreateUser(db, payload, cfg, logger)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
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

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	logger.PrintfInfo("Getting user: %s", user.(*auth.JWTPayload).UserId)
	userFromDb, err := GetUserById(db, user.(*auth.JWTPayload), logger)

	if err != nil {
		logger.PrintfError("Error getting user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}
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

	user, ok := c.Get("user")
	if !ok {
		logger.PrintfError("User data not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	logger.PrintfInfo("Getting profile picture for user: %s", user.(*auth.JWTPayload).UserId)
	imageURL, err := GenerateGetProfilePictureURL(db, user.(*auth.JWTPayload), logger, cfg)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

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
	if email == ":email" {
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:  http.StatusBadRequest,
			Error: enum.MalformedRequest,
		})
		return
	}

	logger.PrintfInfo("Checking if user with email '%s' exists", email)
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

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	logger.PrintfInfo("Updating user: %s", user.(*auth.JWTPayload).UserId)
	updatedUser, err := UpdateUser(db, user.(*auth.JWTPayload), payload, logger)

	if err != nil {
		logger.PrintfError("Error updating user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

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

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	logger.PrintfInfo("Generating upload profile picture url for user: %s", user.(*auth.JWTPayload).UserId)
	uploadURL, err := GenerateUploadProfilePictureURL(db, user.(*auth.JWTPayload), logger, cfg)

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

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	logger.PrintfInfo("Deleting user: %s", user.(*auth.JWTPayload).UserId)
	err := DeleteUser(db, user.(*auth.JWTPayload), logger)

	if err != nil {
		logger.PrintfError("Error deleting user: %s", err.Error)
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, gin.H{})
}
