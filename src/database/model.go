package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	Id        string    `gorm:"type:varchar(36);primaryKey"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Content   string    `gorm:"type:text"`
	Iv        string    `gorm:"type:varchar(25)"`
	ChatId    string    `gorm:"type:varchar(36);index"`
	SenderId  string    `gorm:"type:varchar(36);index"`
	Chat      Chat      `gorm:"foreignKey:ChatId"`
	Sender    User      `gorm:"foreignKey:SenderId"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = uuid.NewString()
	return
}

type Chat struct {
	Id          string    `gorm:"type:varchar(36);primaryKey"`
	CreatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Name        string    `gorm:"type:varchar(255)"`
	Picture     *string   `gorm:"type:varchar(2048)"` // TODO: adjust for s3 file key
	Description *string   `gorm:"type:text"`
	Messages    []Message `gorm:"foreignKey:ChatId"`
}

func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	c.Id = uuid.NewString()
	return
}

type User struct {
	Id         string         `gorm:"type:varchar(36);primaryKey"`
	CreatedAt  time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Email      string         `gorm:"type:varchar(255);uniqueIndex"`
	Password   string         `gorm:"type:text"`
	Name       string         `gorm:"type:varchar(50)"`
	Bio        *string        `gorm:"type:varchar(1000)"`
	Iv         string         `gorm:"type:varchar(25)"`
	PublicKey  string         `gorm:"type:text"`
	PrivateKey string         `gorm:"type:text"`
	Keys       []ChatUserKeys `gorm:"foreignKey:UserId"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.NewString()
	return
}

type ChatUserKeys struct {
	Id        string    `gorm:"type:varchar(36);primaryKey"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;column:updated_at"`
	Key       string    `gorm:"type:text"`
	ChatId    string    `gorm:"type:varchar(36);index"`
	Chat      Chat      `gorm:"foreignKey:ChatId"`
	UserId    string    `gorm:"type:varchar(36);index"`
	User      User      `gorm:"foreignKey:UserId"`
}

func (cuk *ChatUserKeys) BeforeCreate(tx *gorm.DB) (err error) {
	cuk.Id = uuid.NewString()
	return
}

type UserKeys struct {
	Id           string    `gorm:"type:varchar(36);primaryKey"`
	CreatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:datetime"`
	ExpiredAt    time.Time `gorm:"type:datetime"`
	User         User      `gorm:"foreignKey:UserId"`
	UserId       string    `gorm:"type:varchar(36);index"`
	RefreshToken string    `gorm:"type:text"`
}

// doesn't need a before create cause the uuid is determinde by the unique factor of the refresh token
