package chat

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/database"
	"easyflow-backend/src/enum"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func CreateChat(db *gorm.DB, payload *CreateChatRequest, jwtPayload *auth.JWTPayload) (*CreateChatResponse, *api.ApiError) {
	var users []database.User
	var userKeys []UserKeyEntry

	// Start a transaction
	tx := db.Begin()

	//get usres from payload.UserKeys
	for _, userKey := range payload.UserKeys {
		user := database.User{}
		if err := tx.Where("id = ?", userKey.UserID).First(&user).Error; err != nil {
			tx.Rollback()
			return nil, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}
		users = append(users, user)
		userKeys = append(userKeys, userKey)
	}

	if len(users) != len(payload.UserKeys) {
		tx.Rollback()
		return nil, &api.ApiError{
			Code:  http.StatusNotFound,
			Error: enum.UserNotFound,
		}
	}

	if len(users) != len(userKeys) {
		tx.Rollback()
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	chat := &database.Chat{
		Name:        payload.Name,
		Picture:     payload.Picture,
		Description: payload.Description,
		Messages:    nil,
	}

	if err := tx.Create(chat).Error; err != nil {
		tx.Rollback()
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	for i, user := range users {
		chatUserKeys := &database.ChatUserKeys{
			ChatId: chat.Id,
			UserId: user.Id,
			Key:    userKeys[i].Key,
		}

		if err := tx.Create(chatUserKeys).Error; err != nil {
			tx.Rollback()
			return nil, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	return &CreateChatResponse{
		Id:          chat.Id,
		CreatedAt:   chat.CreatedAt.String(),
		UpdateAt:    chat.UpdatedAt.String(),
		Name:        chat.Name,
		Picture:     chat.Picture,
		Description: chat.Description,
	}, nil
}

func GetChatPreviews(db *gorm.DB, jwtPayload *auth.JWTPayload) ([]GetChatPreviewResponse, *api.ApiError) {
	fmt.Println("Attempting to get chat previews for user: ", jwtPayload.UserId)
	var chatUserKeys []database.ChatUserKeys
	var chatPreviews []GetChatPreviewResponse

	if err := db.Where("user_id = ?", jwtPayload.UserId).Find(&chatUserKeys).Error; err != nil {
		fmt.Println("Error getting chats for user: ", err)
		return nil, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	for _, chatUserKey := range chatUserKeys {
		var chat database.Chat
		if err := db.Where("id = ?", chatUserKey.ChatId).First(&chat).Error; err != nil {
			fmt.Println("Error getting chat with id:", chatUserKey.ChatId, ". Error:", err)
			return nil, &api.ApiError{
				Code:  http.StatusInternalServerError,
				Error: enum.ApiError,
			}
		}

		var lastMessage *database.Message = nil
		if err := db.Where("chat_id = ?", chatUserKey.ChatId).Order("created_at desc").First(&lastMessage).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				// If there's another error, log it and return
				fmt.Println("Error getting last message for chat with id: ", chatUserKey.ChatId, ". Error: ", err.Error())
				return nil, &api.ApiError{
					Code:  http.StatusInternalServerError,
					Error: enum.ApiError,
				}
			}

			chatPreview := GetChatPreviewResponse{
				CreateChatResponse: CreateChatResponse{
					Id:          chat.Id,
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

	return chatPreviews, nil
}
