package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func AddCros() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			//AllowAllOrigins:  true,
			AllowOrigins:     []string{"*"}, // 等同于允许所有域名 #AllowAllOrigins:  true
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"*", "Authorization"},
			ExposeHeaders:    []string{"Content-Length", "text/plain", "Authorization", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	)
	//return func(c *gin.Context) {
	//	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	//	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
	//	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	//
	//	if c.Request.Method == "OPTIONS" {
	//		c.AbortWithStatus(http.StatusNoContent)
	//		return
	//	}
	//
	//	c.Next()
	//}
}
