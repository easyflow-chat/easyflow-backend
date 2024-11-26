package user

import (
	"easyflow-backend/api"
	"easyflow-backend/api/auth"
	"easyflow-backend/common"
	"easyflow-backend/enum"
	"easyflow-backend/middleware"
	"net/http"

	"github.com/easyflow-chat/easyflow-backend/lib/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("User"))
	r.Use(middleware.NewRateLimiter(1, 2))
	r.POST("/signup", middleware.NewRateLimiter(1, 0), CreateUserController)
	r.GET("/", auth.AuthGuard(), GetUserController)
	r.GET("/exists/:email", UserExists)
	r.GET("/profile-picture", auth.AuthGuard(), GetProfilePictureController)
	r.GET("/upload-profile-picture", auth.AuthGuard(), GenerateUploadProfilePictureURLController)
	r.PUT("/", auth.AuthGuard(), UpdateUserController)
	r.DELETE("/", auth.AuthGuard(), DeleteUserController)
}

func CreateUserController(c *gin.Context) {
	payload, logger, db, cfg, errors := common.SetupEndpoint[CreateUserRequest](c)
	if len(errors) > 0 {
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
	if len(errors) > 0 {
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

	userFromDb, err := GetUserById(db, user.(*jwt.JWTTokenPayload), logger)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(200, userFromDb)
}

func GetProfilePictureController(c *gin.Context) {
	_, logger, db, cfg, errors := common.SetupEndpoint[any](c)
	if len(errors) > 0 {
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

	imageURL, err := GenerateGetProfilePictureURL(db, user.(*jwt.JWTTokenPayload), logger, cfg)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, imageURL)
}

func UserExists(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[any](c)
	if len(errors) > 0 {
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

	userInDb, err := GetUserByEmail(db, email, logger)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, userInDb)
}

func UpdateUserController(c *gin.Context) {
	payload, logger, db, _, errors := common.SetupEndpoint[UpdateUserRequest](c)
	if len(errors) > 0 {
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

	updatedUser, err := UpdateUser(db, user.(*jwt.JWTTokenPayload), payload, logger)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, updatedUser)
}

func GenerateUploadProfilePictureURLController(c *gin.Context) {
	_, logger, db, cfg, errors := common.SetupEndpoint[any](c)
	if len(errors) > 0 {
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

	uploadURL, err := GenerateUploadProfilePictureURL(db, user.(*jwt.JWTTokenPayload), logger, cfg)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, uploadURL)
}

func DeleteUserController(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[CreateUserRequest](c)
	if len(errors) > 0 {
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

	err := DeleteUser(db, user.(*jwt.JWTTokenPayload), logger)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, gin.H{})
}
