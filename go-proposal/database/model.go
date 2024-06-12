package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	Id        string    `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Content   string    `gorm:"type:text"`
	Iv        string    `gorm:"type:varchar(25)"`
	ChatId    string    `gorm:"index"`
	SenderId  string    `gorm:"index"`
	Chat      Chat      `gorm:"foreignKey:ChatId"`
	Sender    User      `gorm:"foreignKey:SenderId"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = uuid.NewString()
	return
}

type Chat struct {
	Id          string    `gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Name        string    `gorm:"type:varchar(255)"`
	Picture     *string   `gorm:"type:varchar(2048)"`
	Description *string   `gorm:"type:text"`
	Messages    []Message `gorm:"foreignKey:ChatId"`
	Users       []User    `gorm:"many2many:user_chats;"`
}

func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	c.Id = uuid.NewString()
	return
}

type User struct {
	Id             string         `gorm:"type:uuid;primaryKey"`
	CreatedAt      time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Email          string         `gorm:"type:varchar(255);unique_index"`
	Password       string         `gorm:"type:text"`
	Name           string         `gorm:"type:varchar(255)"`
	ProfilePicture *string        `gorm:"column:profile_picture;type:varchar(2048)"`
	Bio            *string        `gorm:"type:varchar(1000)"`
	Iv             string         `gorm:"type:varchar(25)"`
	PublicKey      string         `gorm:"type:text"`
	PrivateKey     string         `gorm:"type:text"`
	Chats          []Chat         `gorm:"many2many:user_chats;"`
	Keys           []ChatUserKeys `gorm:"foreignKey:UserId"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.NewString()
	return
}

type ChatUserKeys struct {
	Id        string     `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	Key       string     `gorm:"type:text"`
	ChatId    string     `gorm:"index"`
	Chat      Chat       `gorm:"foreignKey:ChatId"`
	UserId    string     `gorm:"index"`
	User      User       `gorm:"foreignKey:UserId"`
}

func (cuk *ChatUserKeys) BeforeCreate(tx *gorm.DB) (err error) {
	cuk.Id = uuid.NewString()
	return
}
