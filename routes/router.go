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
	// 中间件
	r.Use(middleware.AddCros())

	r.POST("/createUser", api.AddUser)
	r.DELETE("/deleteUser/:id", api.DeleteUser)
	r.PUT("/updateUser/:id", api.UpdateUser)
	r.GET("/getUser/:id", api.GetUser)
	r.GET("/getUserPage", api.GetUserPage)

	r.POST("/login", api.Login)
	auth := r.Group("api")
	auth.Use(middleware.JwtToken())
	{
		auth.POST("/register", api.AddUser)
	}
	// auth.Use(middleware.JWTAuth())
	// {
	// 	auth.GET("/users", api.GetUsers)
	// }

	r.Run(utils.HttpPort)
}
