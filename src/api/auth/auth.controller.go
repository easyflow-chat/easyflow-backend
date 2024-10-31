package auth

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func extractDomain(c *gin.Context, urlString string) *string {
	origin, err := url.Parse(c.GetHeader("Origin"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err.Error(),
		})
		return nil
	}
	domain := origin.Hostname()
	return &domain
}

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

	tokens, err := LoginService(db, cfg, payload, logger)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	domain := extractDomain(c, c.GetHeader("Origin"))

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokens.AccessToken, cfg.JwtExpirationTime, "/", *domain, cfg.Stage == "production", true)
	c.SetCookie("refresh_token", tokens.RefreshToken, cfg.RefreshExpirationTime, "/", *domain, cfg.Stage == "production", true)

	c.JSON(200, gin.H{})
}

func CheckLoginController(c *gin.Context) {
	_, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: "User not found in context",
		})
		return
	}

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
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	tokens, err := RefreshService(db, cfg, payload.(*JWTAccessTokenPayload), logger)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	domain := extractDomain(c, c.GetHeader("Origin"))

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokens.AccessToken, cfg.JwtExpirationTime, "/", *domain, cfg.Stage == "production", true)
	c.SetCookie("refresh_token", tokens.RefreshToken, cfg.RefreshExpirationTime, "/", *domain, cfg.Stage == "production", true)

	c.JSON(200, gin.H{})
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
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    http.StatusBadRequest,
			Error:   enum.InvalidRefresh,
			Details: err,
		})
	}

	payload, err := ValidateToken(cfg, refresh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		})
		return
	}

	e := LogoutService(db, payload, logger)
	if e != nil {
		c.JSON(e.Code, e)
		return
	}

	domain := extractDomain(c, c.GetHeader("Origin"))

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", 0, "/", *domain, cfg.Stage == "production", true)
	c.SetCookie("refresh_token", "", 0, "/", *domain, cfg.Stage == "production", true)

	c.JSON(200, gin.H{})
}
