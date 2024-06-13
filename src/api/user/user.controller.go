package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.POST("/signup", CreateUserController)
}

func CreateUserController(c *gin.Context) {
	var payload CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, ok := c.Get("db")
	if !ok {
		c.JSON(500, gin.H{"error": "database not found"})
		return
	}

	user, err := CreateUser(db.(*gorm.DB), &payload)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, user)
}
