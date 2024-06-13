package chat

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/enum"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthEndpoints(r *gin.RouterGroup) {
	r.Use(auth.AuthGuard())
	r.POST("", CreateChatController)
}

func CreateChatController(c *gin.Context) {
	var payload CreateChatRequest
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

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	jwtPayload, ok := user.(auth.JWTPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	chat, err := CreateChat(db.(*gorm.DB), &payload, &jwtPayload)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusCreated, chat)
}
