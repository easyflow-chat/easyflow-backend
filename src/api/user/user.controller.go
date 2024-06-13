package user

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.POST("/signup", CreateUserController)
}

func CreateUserController(c *gin.Context) {
	var payload CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:     http.StatusBadRequest,
			Message:  "Failed to parse request body",
			Expected: common.StringPointer("email, name, password, public_key, private_key, iv"),
		})
		return
	}

	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ErrDatabaseConnection)
		return
	}

	user, err := CreateUser(db.(*gorm.DB), &payload)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, user)
}
