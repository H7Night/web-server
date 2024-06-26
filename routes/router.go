package routes

import (
	"github.com/gin-gonic/gin"
	"web-server/controllers"
	"web-server/middleware"
	"web-server/utils"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.AddCros()) // Add the CORS middleware

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
