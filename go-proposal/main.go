package main

import (
	"easflow-backend/database"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Use the actual DSN for your MySQL connection here
	//dsn := "mysql://devel:devel@localhost:3306/chat-app?charset=utf8mb4&parseTime=True&loc=Local"
	//read from env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Use an environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	log.Println("Database URL:", databaseURL)
	dbInst, err := database.NewDatabaseInst(databaseURL, &gorm.Config{})
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
		Chats:          nil,
	})

	defer dbInst.Release() // En	sure to release when done
}
