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
)

func main() {
	cfg := common.LoadDefaultConfig()
	dbInst, err := database.NewDatabaseInst(cfg.DatabaseURL, &cfg.GormConfig)
	if err != nil {
		panic(err)
	}

	logger := common.NewLogger(os.Stdout, "Main")

	err = dbInst.Migrate()
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(cors.Default())
	router.Use(middleware.DatabaseMiddleware(dbInst.GetClient()))
	router.Use(middleware.ConfigMiddleware(cfg))
	router.Use(GinLoggerMiddleware(logger))
	router.Use(gin.Recovery())

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

func GinLoggerMiddleware(logger *common.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		client := c.ClientIP()
		method := c.Request.Method
		latency := time.Now().Sub(start)
		status := c.Writer.Status()
		c.Next()

		logger.Printf("%d |\t%s|\t%s|\t%s %s",
			status,
			latency,
			client,
			method,
			path,
		)
	}
}
