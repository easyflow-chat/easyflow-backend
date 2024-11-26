package chat

import (
	"easyflow-backend/api"
	"easyflow-backend/api/auth"
	"easyflow-backend/common"
	"easyflow-backend/enum"
	"easyflow-backend/middleware"
	"net/http"

	"github.com/easyflow-chat/easyflow-backend/lib/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterChatEndpoints(r *gin.RouterGroup) {
	r.Use(middleware.LoggerMiddleware("Chat"))
	r.Use(auth.AuthGuard())
	r.Use(middleware.NewRateLimiter(1, 5))
	r.POST("/", CreateChatController)
	r.GET("/preview", GetChatPreviewsController)
	r.GET("/:chatId", GetChatByIdController)
}

func CreateChatController(c *gin.Context) {
	payload, logger, db, _, errors := common.SetupEndpoint[CreateChatRequest](c)
	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	logger.PrintfDebug("Payload: %s", payload.Name)

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	chat, err := CreateChat(db, payload, user.(*jwt.JWTTokenPayload), logger)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusCreated, chat)
}

func GetChatPreviewsController(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[any](c)
	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	chats, err := GetChatPreviews(db, user.(*jwt.JWTTokenPayload), logger)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, chats)
}

func GetChatByIdController(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[any](c)
	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	chatId := c.Param("chatId")

	chat, err := GetChatById(db, chatId, user.(*jwt.JWTTokenPayload), logger)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, chat)
}
