package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseInst struct {
	client *gorm.DB
}

func NewDatabaseInst(url string, config *gorm.Config) (*DatabaseInst, error) {
	db, err := gorm.Open(mysql.Open(url), config)
	if err != nil {
		return nil, err
	}

	return &DatabaseInst{client: db}, nil
}

func (d *DatabaseInst) GetClient() *gorm.DB {
	return d.client
}

func (d *DatabaseInst) Migrate() error {
	return d.client.AutoMigrate(&Message{}, &Chat{}, &User{}, &ChatUserKeys{}, &UserKeys{})
}

func GetUserByEmail(db *gorm.DB, email *string) (*User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
