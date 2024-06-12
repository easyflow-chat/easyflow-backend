package main

import (
	"easflow-backend/config"
	"easflow-backend/database"

	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadDefaultConfig()
	dbInst, err := database.NewDatabaseInst(cfg.DatabaseURL, &gorm.Config{})
	if err != nil {
		panic(err) // Handle error properly in production code
	}
	dbInst.Acquire()
	dbInst.Migrate()
	dbInst.GetClient().Create(&database.User{
		Email:          "test@test.com",
		Password:       "password",
		Name:           "Test User",
		ProfilePicture: nil,
		Bio:            nil,
		Iv:             "123456789",
		PublicKey:      "publickey",
		PrivateKey:     "privatekey",
		Chats:          nil,
	})

	defer dbInst.Release() // En	sure to release when done
}
