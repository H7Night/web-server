package routes

import (
	api "web-server/api"
	"web-server/middleware"
	"web-server/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.AddCros())

	r.POST("/adduser", api.AddUser)
	r.DELETE("/deleteuser", api.DeleteUser)

	auth := r.Group("api")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/register", api.AddUser)
	}
	// auth.Use(middleware.JWTAuth())
	// {
	// 	auth.GET("/users", api.GetUsers)
	// }

	// public := r.Group("api")
	// {
	// 	public.POST("/login", api.Login)
	// }

	r.Run(utils.HttpPort)
}
