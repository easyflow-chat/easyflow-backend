package main

import (
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

	dbInst.Acquire()
	err = dbInst.Migrate()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middleware.DatabaseMiddleware(dbInst.GetClient()))

	//register user endpoints
	userEndpoints := router.Group("/user")
	{
		user.RegisterUserEndpoints(userEndpoints)
	}

	router.Run(":" + cfg.Port)

	defer dbInst.Release()
}
