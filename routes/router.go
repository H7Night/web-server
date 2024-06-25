package routes

import (
	"github.com/gin-gonic/gin"
	"goIland/controllers"
	"goIland/middleware"
	"goIland/utils"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	auth := r.Group("api/v1")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/users", controllers.GetUsers)
	}

	public := r.Group("api/v1")
	{
		public.POST("/login", controllers.Login)
	}

	r.Run(utils.HttpPort)
}
