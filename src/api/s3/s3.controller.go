package s3

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterS3Endpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("S3"))
	r.Use(middleware.RateLimiter(1, 4))
	r.POST("/list-objects", GetObjectsController)
	r.POST("/get-download-url", GetDownloadURLController)
}

func GetObjectsController(c *gin.Context) {
	payload, logger, _, cfg, errors := common.SetupEndpoint[BucketRequest](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	objects, err := GetObjects(logger, cfg, payload.Name)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	for _, object := range objects.Contents {
		logger.PrintfInfo("Object in bucket: %s", *object.Key)
	}

	c.JSON(200, "Successfully fetched objects")
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
