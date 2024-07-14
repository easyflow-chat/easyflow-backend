package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAuthEndpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("Auth"))
	r.Use(middleware.RateLimiter(1, 0))
	r.POST("/login", LoginController)
	r.GET("/refresh", RefreshAuthGuard(), RefreshController)
}

func LoginController(c *gin.Context) {
	payload, logger, db, cfg, errors := common.SetupEndpoint[LoginRequest](c, "Auth")
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	logger.Printf("Successfully validated request for login")

	jwtPair, err := LoginService(db, cfg, payload, logger)
	if err != nil {
		logger.PrintfError("Error logging in: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	c.SetCookie("access_token", jwtPair.AccessToken, TOK_TTL, "/", "", false, true)
	c.SetCookie("refresh_token", jwtPair.RefreshToken, TOK_TTL, "/", "", false, true)

	logger.Println("Responding with 200 status code for successful login")
	c.JSON(200, &gin.H{})
}

func RefreshController(c *gin.Context) {
	_, logger, db, cfg, errors := common.SetupEndpoint[any](c, "Auth")
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	logger.Printf("Successfully validated request for token refresh")

	payload, ok := c.Get("user")
	if !ok {
		logger.PrintfError("User not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	jwtPair, err := RefreshService(db, cfg, payload.(*JWTPayload), logger)

	if err != nil {
		logger.PrintfError("Error refreshing token: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	c.SetCookie("access_token", jwtPair.AccessToken, TOK_TTL, "/", "", false, true)

	logger.Println("Responding with 200 status code for successful token refresh")
	c.JSON(200, &gin.H{})
}
