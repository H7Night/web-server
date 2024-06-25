package routes

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("front", "web/front/dist/index.html")
	return p
}

func InitRouter() {
	r := gin.New()
	// 设置信任网络 []string
	// nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)
	r.HTMLRender = createRender()
	r.Use(gin.Recovery())

	router := r.Group("api/v1")
	{
		router.GET("user/add")
	}
}
