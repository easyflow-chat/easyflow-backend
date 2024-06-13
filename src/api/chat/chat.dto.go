package chat

type CreateChatRequest struct {
	Name        string            `json:"name" validate:"required"`
	Picture     *string           `json:"picture" validate:"omitempty,url"`
	Description *string           `json:"description" validate:"omitempty"`
	Users       []string          `json:"users" validate:"required"`
	UserKeys    map[string]string `json:"user_keys" validate:"required"`
}
