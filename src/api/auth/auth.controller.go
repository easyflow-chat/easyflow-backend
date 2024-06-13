package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/enum"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthEndpoints(r *gin.RouterGroup) {
	r.POST("/login", LoginController)
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

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	jwtPair, err := LoginService(db.(*gorm.DB), &payload)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.SetCookie("access_token", jwtPair.AccessToken, TOK_TTL, "/", "", false, true)
	c.SetCookie("refresh_token", jwtPair.RefreshToken, TOK_TTL, "/", "", false, true)

	c.JSON(200, gin.H{})
}
