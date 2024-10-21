package main

import (
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/api/chat"
	"easyflow-backend/src/api/user"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/middleware"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := common.LoadDefaultConfig()

	log := common.NewLogger(os.Stdout, "Main", nil)
	var isConnected = false
	var dbInst *database.DatabaseInst
	var connectionAttempts = 0
	var connectionPause = 5
	for !isConnected {
		var err error
		dbInst, err = database.NewDatabaseInst(cfg.DatabaseURL, &cfg.GormConfig)

		if err != nil {
			if connectionAttempts <= 5 {
				connectionAttempts++
				log.PrintfError("Failed to connect to database, retrying in %d seconds. Attempt %d", connectionPause, connectionAttempts)
				time.Sleep(time.Duration(connectionPause) * time.Second)
				connectionPause += 5
			} else {
				panic(err)
			}
		} else {
			isConnected = true
		}
	}

	if !cfg.DebugMode {
		gin.SetMode(gin.ReleaseMode)
		dbInst.SetLogMode(logger.Silent)
	}

	err := dbInst.Migrate()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	log.Printf("Frontend URL for cors: %s", cfg.FrontendURL)
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowWildcard:    true,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
	}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.DatabaseMiddleware(dbInst.GetClient()))
	router.Use(middleware.ConfigMiddleware(cfg))
	router.Use(gin.Recovery())

	//register user endpoints
	userEndpoints := router.Group("/user")
	{
		log.Printf("Registering user endpoints")
		user.RegisterUserEndpoints(userEndpoints)
	}

	authEndpoints := router.Group("/auth")
	{
		log.Printf("Registering auth endpoints")
		auth.RegisterAuthEndpoints(authEndpoints)
	}

	chatEndpoints := router.Group("/chat")
	{
		log.Printf("Registering chat endpoints")
		chat.RegisterChatEndpoints(chatEndpoints)
	}

	log.Printf("Starting server on port %s", cfg.Port)
	router.Run(":" + cfg.Port)
}
