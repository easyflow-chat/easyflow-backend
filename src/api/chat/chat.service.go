package chat

import (
	"easyflow-backend/src/api"
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/database"

	"gorm.io/gorm"
)

func CreateChat(db *gorm.DB, payload *CreateChatRequest, jwtPayload *auth.JWTPayload) (*database.Chat, *api.ApiError) {
	return nil, nil
}
