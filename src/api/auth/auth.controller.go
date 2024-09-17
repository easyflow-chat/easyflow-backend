package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"net/http"
	"time"

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

	logger.PrintfInfo("Logging in user with email: %s", payload.Email)
	jwtPair, err := LoginService(db, cfg, payload, logger)
	if err != nil {
		logger.PrintfError("Error logging in: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", jwtPair.AccessToken, cfg.RefreshExpirationTime, "/", cfg.BackendDomain, cfg.Stage != "development", false)
	c.SetCookie("refresh_token", jwtPair.RefreshToken, cfg.RefreshExpirationTime, "/", cfg.BackendDomain, cfg.Stage != "development", false)

	c.JSON(200, true)
}

func CheckLoginController(c *gin.Context) {
	_, logger, _, _, errors := common.SetupEndpoint[any](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		logger.PrintfError("User not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: "User not found in context",
		})
		return
	}

	logger.PrintfInfo("User with id: %s is logged in", user.(*JWTPayload).UserId)
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

	logger.PrintfInfo("Refreshing token for user with id: %s", payload.(*JWTPayload).UserId)
	accessToken, err := RefreshService(db, cfg, payload.(*JWTPayload), logger)

	if err != nil {
		logger.PrintfError("Error refreshing token: %s", err.Details)
		// in case of error, clear the cookies so the user has to log in again
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("access_token", "", -1, "/", cfg.BackendDomain, cfg.Stage != "development", false)
		c.SetCookie("refresh_token", "", -1, "/", cfg.BackendDomain, cfg.Stage != "development", false)
		c.JSON(err.Code, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", *accessToken, int(payload.(*JWTPayload).ExpiresAt.Unix())-int(time.Now().Unix()), "/", cfg.BackendDomain, cfg.Stage != "development", false)

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
		logger.PrintfWarning("No refresh token")
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

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", cfg.BackendDomain, cfg.Stage != "development", false)
	c.SetCookie("refresh_token", "", -1, "/", cfg.BackendDomain, cfg.Stage != "development", false)

	c.JSON(200, gin.H{})
	logger.Printf("Successfully logged out user with id: %s", payload.UserId)
}
