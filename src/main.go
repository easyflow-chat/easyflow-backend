package main

import (
	"easyflow-backend/src/api/auth"
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

	testAuthEndpoints := router.Group("/test")
	{
		testAuthEndpoints.Use(auth.AuthGuard())
		func(r *gin.RouterGroup) {
			r.GET("/auth", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "success",
				})
			})
		}(testAuthEndpoints)
	}

	router.Run(":" + cfg.Port)
}
