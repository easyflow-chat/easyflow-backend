package chat

type UserKeyEntry struct {
	UserID string `json:"userId" validate:"required"`
	Key    string `json:"key" validate:"required"`
}

type CreateChatRequest struct {
	Name        string         `json:"name" validate:"required"`
	Picture     *string        `json:"picture" validate:"omitempty,url"`
	Description *string        `json:"description" validate:"omitempty"`
	UserKeys    []UserKeyEntry `json:"userKeys" validate:"required"`
}

type CreateChatResponse struct {
	Id          string  `json:"id"`
	CreatedAt   string  `json:"created_at"`
	UpdateAt    string  `json:"updated_at"`
	Name        string  `json:"name"`
	Picture     *string `json:"picture"`
	Description *string `json:"description"`
}
