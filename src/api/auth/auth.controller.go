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
	r.Use(middleware.RateLimiter(1, 2))
	r.POST("/login", LoginController)
	r.GET("/check", AuthGuard(), CheckLoginController)
	r.GET("/refresh", RefreshAuthGuard(), RefreshController)
	r.GET("/logout", AuthGuard(), LogoutController)
}

func LoginController(c *gin.Context) {
	payload, logger, db, cfg, errors := common.SetupEndpoint[LoginRequest](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	jwtPair, err := LoginService(db, cfg, payload, logger)
	if err != nil {
		logger.PrintfError("Error logging in: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	c.SetCookie("access_token", jwtPair.AccessToken, cfg.JwtExpirationTime, "/", "", cfg.Stage != "development", true)
	// make it valid for 10 more minutes than the access token so it can be refreshed
	c.SetCookie("refresh_token", jwtPair.RefreshToken, cfg.RefreshExpirationTime, "/", "", cfg.Stage != "development", true)

	c.JSON(200, true)
}

func CheckLoginController(c *gin.Context) {
	// only returns if it comes through the authguard so we can assume the user is logged in
	c.JSON(200, true)
}

func RefreshController(c *gin.Context) {
	_, logger, db, cfg, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

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
		// in case of error, clear the cookies so the user has to log in again
		c.SetCookie("access_token", "", -1, "/", "", cfg.Stage != "development", true)
		c.SetCookie("refresh_token", "", -1, "/", "", cfg.Stage != "development", true)
		c.JSON(err.Code, err)
		return
	}

	c.SetCookie("access_token", jwtPair.AccessToken, cfg.JwtExpirationTime, "/", "", cfg.Stage != "development", true)

	c.JSON(200, &gin.H{})
}

func LogoutController(c *gin.Context) {
	_, logger, db, cfg, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		logger.PrintfError("No refresh token")
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    http.StatusBadRequest,
			Error:   enum.InvalidRefresh,
			Details: err,
		})
	}

	payload, err := ValidateToken(cfg, refresh)
	if err != nil {
		logger.PrintfError("An error occoured while validating refresh token")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		})
		return
	}

	logger.PrintfInfo("Trying to logout user with id: %s", payload.UserId)
	e := LogoutService(db, payload, logger)
	if e != nil {
		logger.PrintfError("An error occured while logging out user with id: %s", payload.UserId)
		c.JSON(e.Code, e)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", cfg.Stage != "development", true)
	c.SetCookie("refresh_token", "", -1, "/", "", cfg.Stage != "development", true)

	c.JSON(200, gin.H{})
	logger.Printf("Successfully logged out user with id: %s", payload.UserId)
}
