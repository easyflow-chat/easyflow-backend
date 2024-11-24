package chat

import (
	"easyflow-backend/api"
	"easyflow-backend/enum"
	"net/http"

	"github.com/easyflow-chat/easyflow-backend/lib/database"
	"github.com/easyflow-chat/easyflow-backend/lib/jwt"
	"github.com/easyflow-chat/easyflow-backend/lib/logger"

	"gorm.io/gorm"
)

func CreateChat(db *gorm.DB, payload *CreateChatRequest, jwtPayload *jwt.JWTTokenPayload, logger *logger.Logger) (*CreateChatResponse, *api.ApiError) {
	var users []database.User
	var userKeys []UserKeyEntry

	//get users from payload.UserKeys
	for _, userKey := range payload.UserKeys {
		user := database.User{ID: jwtPayload.UserID}
		if err := db.First(&user).Error; err != nil {
			logger.PrintfError("Error getting user with id: %s", userKey.UserID)
			return nil, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}
		users = append(users, user)
		userKeys = append(userKeys, userKey)
	}

	if len(users) != len(payload.UserKeys) {
		logger.PrintfError("User keys and users length mismatch")
		return nil, &api.ApiError{
			Code:  http.StatusNotFound,
			Error: enum.UserNotFound,
		}
	}

	if len(users) != len(userKeys) {
		logger.PrintfError("User keys and users length mismatch")
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	logger.PrintfDebug("Payload: %s", payload.Name)

	chat := database.Chat{
		Name:        payload.Name,
		Picture:     payload.Picture,
		Description: payload.Description,
	}

	// Start a transaction
	tx := db.Begin()

	if err := tx.Create(&chat).Error; err != nil {
		tx.Rollback()
		logger.PrintfError("Error creating chat: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	for i, user := range users {
		chatUserKeys := database.ChatsUsers{
			ChatID: chat.ID,
			UserID: user.ID,
			Key:    userKeys[i].Key,
		}

		if err := tx.Create(&chatUserKeys).Error; err != nil {
			tx.Rollback()
			logger.PrintfError("Error creating chat user key: %s", err)
			return nil, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.PrintfError("Error committing transaction: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	logger.Printf("Successfully created chat with id: %s", chat.ID)

	return &CreateChatResponse{
		Id:          chat.ID,
		CreatedAt:   chat.CreatedAt.String(),
		UpdateAt:    chat.UpdatedAt.String(),
		Name:        chat.Name,
		Picture:     chat.Picture,
		Description: chat.Description,
	}, nil
}

func GetChatPreviews(db *gorm.DB, jwtPayload *jwt.JWTTokenPayload, logger *logger.Logger) ([]GetChatPreviewResponse, *api.ApiError) {
	logger.PrintfInfo("Attempting to get chat previews for user: %s", jwtPayload.UserID)
	var chatUserKeys []database.ChatsUsers
	chatPreviews := []GetChatPreviewResponse{}

	if err := db.Where("user_id = ?", jwtPayload.UserID).Find(&chatUserKeys).Error; err != nil {
		logger.PrintfError("Error getting chats for user: %s", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	for _, chatUserKey := range chatUserKeys {
		var chat database.Chat
		if err := db.Where("id = ?", chatUserKey.ChatID).First(&chat).Error; err != nil {
			logger.PrintfError("Error getting chat with id: %s. %s", chatUserKey.ChatID, err)
			return nil, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}

		var lastMessage *database.Message = nil
		if err := db.Where("chat_id = ?", chatUserKey.ChatID).Order("created_at desc").First(&lastMessage).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				// If there's another error, log it and return
				logger.PrintfError("Error getting last message for chat with id: %s. Error: %s", chatUserKey.ChatID, err.Error())
				return nil, &api.ApiError{
					Code:  http.StatusInternalServerError,
					Error: enum.ApiError,
				}
			}

			chatPreview := GetChatPreviewResponse{
				CreateChatResponse: CreateChatResponse{
					Id:          chat.ID,
					CreatedAt:   chat.CreatedAt.String(),
					UpdateAt:    chat.UpdatedAt.String(),
					Name:        chat.Name,
					Picture:     chat.Picture,
					Description: chat.Description,
				},
				LastMessage: &lastMessage.Content,
			}

			chatPreviews = append(chatPreviews, chatPreview)
		}
	}

	logger.Printf("Successfully got chat previews for user: %s", jwtPayload.UserID)

	return chatPreviews, nil
}

func GetChatById(db *gorm.DB, chatId string, jwtPayload *jwt.JWTTokenPayload, logger *logger.Logger) (*GetChatByIdResponse, *api.ApiError) {
	var chat database.Chat
	if err := db.Where("id = ?", chatId).First(&chat).Error; err != nil {
		logger.PrintfError("Error getting chat with id: %s. Error: %s", chatId, err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	var chatUserKeys []database.ChatsUsers
	if err := db.Where("chat_id = ? AND user_id = ?", chatId, jwtPayload.UserID).Find(&chatUserKeys).Error; err != nil {
		logger.PrintfError("Error getting chat user key for chat with id: %s. Error: %s", chatId, err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	var Messages []database.Message
	if err := db.Where("chat_id = ?", chatId).Order("created_at desc").Find(&Messages).Limit(50).Error; err != nil {
		logger.PrintfError("Error getting messages for chat with id: %s. Error: %s", chatId, err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	// Mappings
	usersEntries := []UserEntry{}
	for _, chatUserKey := range chatUserKeys {
		var user database.User
		if err := db.Where("id = ?", chatUserKey.UserID).First(&user).Error; err != nil {
			logger.PrintfError("Error getting user with id: %s. Error: %s", chatUserKey.UserID, err)
			return nil, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}

		usersEntries = append(usersEntries,
			UserEntry{
				Id:   user.ID,
				Name: user.Name,
				Bio:  user.Bio,
			},
		)
	}

	// TODO: Just make one object for user keys not array
	userKeyEntries := []UserKeyEntry{}
	for _, chatUserKey := range chatUserKeys {
		userKeyEntries = append(userKeyEntries,
			UserKeyEntry{
				UserID: chatUserKey.UserID,
				Key:    chatUserKey.Key,
			},
		)
	}

	messageEntries := []MessageEntry{}
	for _, message := range Messages {
		messageEntries = append(messageEntries,
			MessageEntry{
				Id:        message.ID,
				CreatedAt: message.CreatedAt.String(),
				UpdatedAt: message.UpdatedAt.String(),
				Content:   message.Content,
				Iv:        message.Iv,
				SenderId:  message.SenderID,
			},
		)
	}

	logger.Printf("Successfully got chat with id: %s", chatId)

	return &GetChatByIdResponse{
		CreateChatResponse: CreateChatResponse{
			Id:          chat.ID,
			CreatedAt:   chat.CreatedAt.String(),
			UpdateAt:    chat.UpdatedAt.String(),
			Name:        chat.Name,
			Picture:     chat.Picture,
			Description: chat.Description,
		},
		Users:    usersEntries,
		UserKeys: userKeyEntries,
		Messages: messageEntries,
	}, nil

}
