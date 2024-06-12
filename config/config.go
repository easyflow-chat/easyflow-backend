package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	//gorm
	GormConfig gorm.Config
	//env
	DatabaseURL string
	SaltRounds  int
	Port        string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func LoadDefaultConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		GormConfig: gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		},
		DatabaseURL: getEnv("DATABASE_URL", ""),
		SaltRounds:  getEnvInt("SALT_OR_ROUNDS", 10),
		Port:        getEnv("PORT", "8080"),
	}
}
