package common

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/easyflow-chat/easyflow-backend/lib/logger"
	gormLogger "gorm.io/gorm/logger"
)

type Config struct {
	//gorm
	GormConfig gorm.Config
	// stage
	Stage string
	// log level
	LogLevel logger.LogLevel
	//env
	DatabaseURL string
	//jwt
	JwtSecret string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// func getEnvInt(key string, fallback int) int {
// 	if value, ok := os.LookupEnv(key); ok {
// 		if i, err := strconv.Atoi(value); err == nil {
// 			return i
// 		}
// 	}
// 	return fallback
// }

// LoadDefaultConfig loads the default configuration values.
// It reads the environment variables from the .env file, if present,
// and returns a Config struct with the loaded values.
func LoadDefaultConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	return &Config{
		GormConfig: gorm.Config{
			Logger:                                   gormLogger.Default.LogMode(gormLogger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
		},
		Stage:       getEnv("STAGE", "development"),
		LogLevel:    logger.LogLevel(getEnv("LOG_LEVEL", "DEBUG")),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JwtSecret:   getEnv("JWT_SECRET", "public_secret"),
	}
}
