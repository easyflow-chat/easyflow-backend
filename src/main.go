package main

import (
	"easyflow-backend/src/api/auth"
	"easyflow-backend/src/api/chat"
	"easyflow-backend/src/api/user"
	"easyflow-backend/src/common"
	"easyflow-backend/src/database"
	"easyflow-backend/src/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := common.LoadDefaultConfig()
	dbInst, err := database.NewDatabaseInst(cfg.DatabaseURL, &cfg.GormConfig)
	if err != nil {
		panic(err)
	}

	err = dbInst.Migrate()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middleware.DatabaseMiddleware(dbInst.GetClient()))
	router.Use(middleware.ConfigMiddleware(cfg))

	//register user endpoints
	userEndpoints := router.Group("/user")
	{
		user.RegisterUserEndpoints(userEndpoints)
	}

	authEndpoints := router.Group("/auth")
	{
		auth.RegisterAuthEndpoints(authEndpoints)
	}

	chatEndpoints := router.Group("/chat")
	{
		chat.RegisterChatEndpoints(chatEndpoints)
	}

	router.Run(":" + cfg.Port)
}
