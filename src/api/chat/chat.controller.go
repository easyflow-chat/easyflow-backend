package chat

import (
	"easyflow-backend/src/api/auth"

	"github.com/gin-gonic/gin"
)

func RegisterAuthEndpoints(r *gin.RouterGroup) {
	r.Use(auth.AuthGuard())
}
