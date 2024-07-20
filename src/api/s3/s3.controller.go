package s3

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

func RegisterS3Endpoints(r *gin.RouterGroup) {
	r.Use(auth.AuthGuard())
	r.Use(middleware.LoggerMiddleware("S3"))
	r.Use(middleware.RateLimiter(1, 4))
	r.GET("/list-objects/:userid", GetObjectsForUserController)
	r.POST("/get-download-url", GetDownloadURLController)
}

func GetObjectsForUserController(c *gin.Context) {
	_, logger, _, cfg, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	userId := c.Param("userid")

	user, ok := c.Get("user")
	if !ok {
		logger.PrintfError("User not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusUnauthorized,
			Error: enum.UserNotFound,
		})
		return
	}

	if user.(*auth.JWTPayload).UserId != userId {
		logger.PrintfError("User with id %s not authorized to fetch objects for user: %s", user.(*auth.JWTPayload).UserId, userId)
		c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:  http.StatusUnauthorized,
			Error: enum.Unauthorized,
		})
		return
	}

	logger.PrintfInfo("Fetching objects for user: %s", userId)

	objects, err := GetObjectsWithPrefix(logger, cfg, "test", fmt.Sprintf("%s/", userId))
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, objects)
}

func GetDownloadURLController(c *gin.Context) {
	payload, logger, _, cfg, errors := common.SetupEndpoint[GetDownloadURLRequest](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	downloadLink, err := GetDownloadURL(logger, cfg, payload.BucketName, payload.ObjectKey)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, downloadLink)
}
