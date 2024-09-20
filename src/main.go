package main

import (
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/api/chat"
	"easyflow-backend/src/api/user"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/middleware"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := common.LoadDefaultConfig()
	dbInst, err := database.NewDatabaseInst(cfg.DatabaseURL, &cfg.GormConfig)

	if err != nil {
		panic(err)
	}

	if !cfg.DebugMode {
		gin.SetMode(gin.ReleaseMode)
		dbInst.SetLogMode(logger.Silent)
	}

	err = dbInst.Migrate()
	if err != nil {
		panic(err)
	}

	logger := common.NewLogger(os.Stdout, "Main", nil)

	router := gin.New()
	logger.Printf("Frontend URL for cors: %s", cfg.FrontendURL)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowWildcard:    true,
	}))
	router.Use(middleware.DatabaseMiddleware(dbInst.GetClient()))
	router.Use(middleware.ConfigMiddleware(cfg))
	router.Use(gin.Recovery())

	//register user endpoints
	userEndpoints := router.Group("/user")
	{
		logger.Printf("Registering user endpoints")
		user.RegisterUserEndpoints(userEndpoints)
	}

	authEndpoints := router.Group("/auth")
	{
		logger.Printf("Registering auth endpoints")
		auth.RegisterAuthEndpoints(authEndpoints)
	}

	chatEndpoints := router.Group("/chat")
	{
		logger.Printf("Registering chat endpoints")
		chat.RegisterChatEndpoints(chatEndpoints)
	}

	logger.Printf("Starting server on port %s", cfg.Port)
	router.Run(":" + cfg.Port)
}
