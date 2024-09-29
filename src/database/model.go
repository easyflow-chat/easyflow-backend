package database

import (
	"time"
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

type Chat struct {
	Id          string    `gorm:"type:varchar(36);primaryKey"`
	CreatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Name        string    `gorm:"type:varchar(255)"`
	Picture     *string   `gorm:"type:varchar(2048)"` // TODO: adjust for s3 file key
	Description *string   `gorm:"type:text"`
	Messages    []Message `gorm:"foreignKey:ChatId"`
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

type ChatUserKeys struct {
	Id        string    `gorm:"type:varchar(36);primaryKey"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Key       string    `gorm:"type:text"`
	ChatId    string    `gorm:"type:varchar(36);index"`
	Chat      Chat      `gorm:"foreignKey:ChatId"`
	UserId    string    `gorm:"type:varchar(36);index"`
	User      User      `gorm:"foreignKey:UserId"`
}

type UserKeys struct {
	Id           string    `gorm:"type:varchar(36);primaryKey;default:UUID()"`
	CreatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	ExpiredAt    time.Time `gorm:"type:datetime"`
	Random       string    `gorm:"type:varchar(36)"`
	User         User      `gorm:"foreignKey:UserId"`
	UserId       string    `gorm:"type:varchar(36);index"`
	RefreshToken string    `gorm:"type:text"`
}
