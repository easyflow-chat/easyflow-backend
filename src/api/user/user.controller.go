package user

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/enum"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.POST("/signup", CreateUserController)
	r.GET("/:id", GetUserController)
	r.PUT("/:id", UpdateUserController)
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

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	user, err := CreateUser(db.(*gorm.DB), &payload)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, user)
}

func GetUserController(c *gin.Context) {
	id := c.Param("id")

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	user, err := GetUserById(db.(*gorm.DB), &id)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, user)
}

func UpdateUserController(c *gin.Context) {
	id := c.Param("id")

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

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
	}

	user, err := UpdateUser(db.(*gorm.DB), &id, &payload)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, user)
}
