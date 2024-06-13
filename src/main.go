package main

import (
	"easyflow-backend/src/config"
	"easyflow-backend/src/database"
	"easyflow-backend/src/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadDefaultConfig()
	dbInst, err := database.NewDatabaseInst(cfg.DatabaseURL, &cfg.GormConfig)
	if err != nil {
		panic(err)
	}

	dbInst.Acquire()
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middleware.DatabaseMiddleware(dbInst.GetClient()))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(":" + cfg.Port)

	defer dbInst.Release()
}
