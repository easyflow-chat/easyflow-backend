package chat

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"easyflow-backend/src/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterChatEndpoints(r *gin.RouterGroup) {
	r.Use(auth.AuthGuard())
	r.Use(middleware.LoggerMiddleware("Chat"))
	r.POST("", CreateChatController)
	r.GET("/preview", GetChatPreviewsController)
	r.GET("/:chatId", GetChatByIdController)
}

func CreateChatController(c *gin.Context) {
	var payload CreateChatRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    http.StatusBadRequest,
			Error:   enum.MalformedRequest,
			Details: err.Error(),
		})
		return
	}

	if err := api.Validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    http.StatusBadRequest,
			Error:   enum.MalformedRequest,
			Details: api.TranslateError(err),
		})
		return
	}

	raw_logger, ok := c.Get("logger")
	if !ok {
		log.Println("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		log.Println("Type assertion to *common.Logger failed")
	}

	logger.Printf("Successfully validated request for creating chat")

	db, ok := c.Get("db")
	if !ok {
		logger.PrintfError("Database not found in context")
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		})
		return
	}

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

	chat, err := CreateChat(db.(*gorm.DB), &payload, jwtPayload, logger)
	if err != nil {
		logger.PrintfError("Error creating chat: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	logger.Printf("Responding with 201 Created for chat creation")
	c.JSON(http.StatusCreated, chat)
}

func GetChatPreviewsController(c *gin.Context) {
	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
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

	raw_logger, ok := c.Get("logger")
	if !ok {
		log.Println("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		log.Println("Type assertion to *common.Logger failed")
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

	chats, err := GetChatPreviews(db.(*gorm.DB), jwtPayload, logger)
	if err != nil {
		logger.PrintfError("Error getting chat previews: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	logger.Printf("Responding with 200 OK for chat previews")
	c.JSON(http.StatusOK, chats)
}

func GetChatByIdController(c *gin.Context) {
	db, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
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

	raw_logger, ok := c.Get("logger")
	if !ok {
		log.Println("Logger not found in context")
	}

	logger, ok := raw_logger.(*common.Logger)
	if !ok {
		log.Println("Type assertion to *common.Logger failed")
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

	chat, err := GetChatById(db.(*gorm.DB), chatId, jwtPayload, logger)
	if err != nil {
		logger.PrintfError("Error getting chat by id: %s", err.Details)
		c.JSON(err.Code, err)
		return
	}

	logger.Printf("Responding with 200 OK for chat by id")
	c.JSON(http.StatusOK, chat)
}
