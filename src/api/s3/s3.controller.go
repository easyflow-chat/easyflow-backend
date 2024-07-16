package s3

import (
	"bytes"
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterS3Endpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("S3"))
	r.Use(middleware.RateLimiter(1, 4))
	r.GET("/test", GetObjectsController)
	r.POST("/test-upload", UploadFileController)
}

func GetObjectsController(c *gin.Context) {
	_, logger, _, cfg, errors := common.SetupEndpoint[any](c, "S3")
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	GetObjects(logger, cfg)

	c.JSON(200, "Successfully fetched objects")
}

func UploadFileController(c *gin.Context) {
	_, logger, _, cfg, errors := common.SetupEndpoint[any](c, "S3")
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		})
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		})
	}

	UploadFile(logger, cfg, "easyflow", "test", buf, header.Filename)
}
