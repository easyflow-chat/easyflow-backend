package main

import (
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/api/chat"
	"easyflow-backend/src/api/user"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/middleware"
	"os"
	"strings"
	"time"

	cors "github.com/OnlyNico43/gin-cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := common.LoadDefaultConfig()

	log := common.NewLogger(os.Stdout, "Main", nil, common.LogLevel(cfg.LogLevel))
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

	err = router.SetTrustedProxies(nil)
	if err != nil {
		log.PrintfError("Could not set trusted proxies list")
		return
	}

	log.Printf("Frontend URL for cors: %s", cfg.FrontendURL)

	router.Use(cors.CorsMiddleware(cors.Config{
		AllowedOrigins:   strings.Split(cfg.FrontendURL, ", "),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.PrintfError("Failed to start server: %s", err)
		return
	}
}
