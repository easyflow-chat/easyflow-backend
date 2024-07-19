package chat

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterChatEndpoints(r *gin.RouterGroup) {
	r.Use(auth.AuthGuard())
	r.Use(middleware.LoggerMiddleware("Chat"))
	r.Use(middleware.RateLimiter(1, 5))
	r.POST("", CreateChatController)
	r.GET("/preview", GetChatPreviewsController)
	r.GET("/:chatId", GetChatByIdController)
}

func CreateChatController(c *gin.Context) {
	payload, logger, db, _, errors := common.SetupEndpoint[CreateChatRequest](c)
	if errors != nil {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: errors,
		})
		return
	}

	logger.Printf("Successfully validated request for creating chat")

	user, ok := c.Get("user")
	if !ok {
		logger.PrintfError("User not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	jwtPayload, ok := user.(*auth.JWTPayload)
	if !ok {
		logger.PrintfError("Type assertion to *auth.JWTPayload failed")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	chat, err := CreateChat(db, payload, jwtPayload, logger)
	if err != nil {
		logger.PrintfError("Error creating chat: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	logger.Printf("Responding with 201 Created for chat creation")
	c.JSON(http.StatusCreated, chat)
}

func GetChatPreviewsController(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[any](c)
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
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	logger.Printf("Successfully validated request for getting chat previews")

	jwtPayload, ok := user.(*auth.JWTPayload)
	if !ok {
		logger.PrintfError("Type assertion to *auth.JWTPayload failed")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	chats, err := GetChatPreviews(db, jwtPayload, logger)
	if err != nil {
		logger.PrintfError("Error getting chat previews: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	logger.Printf("Responding with 200 OK for chat previews")
	c.JSON(http.StatusOK, chats)
}

func GetChatByIdController(c *gin.Context) {
	_, logger, db, _, errors := common.SetupEndpoint[any](c)
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
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	logger.Printf("Successfully validated request for getting chat by id")

	jwtPayload, ok := user.(*auth.JWTPayload)
	if !ok {
		logger.PrintfError("Type assertion to *auth.JWTPayload failed")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

	chatId := c.Param("chatId")

	chat, err := GetChatById(db, chatId, jwtPayload, logger)
	if err != nil {
		logger.PrintfError("Error getting chat by id: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	logger.Printf("Responding with 200 OK for chat by id")
	c.JSON(http.StatusOK, chat)
}
