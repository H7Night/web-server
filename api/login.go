package api

import (
	"net/http"
	"time"
	"web-server/middleware"
	"web-server/models"
	"web-server/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Login 登录，
func Login(c *gin.Context) {
	var formData models.User
	_ = c.ShouldBind(&formData)
	var code int

	// 从请求中获取 state 参数，用于区分前台登录还是后台登录
	state := c.Query("state")
	if state == "" {
		state = c.PostForm("state")
	}
	//登录校验
	formData, code = models.CheckLogin(formData.Name, formData.Password, state)
	if code == errmsg.Success {
		setToken(c, formData)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"id":      formData.ID,
		"state":   state,
		"data":    formData.Name,
		"message": errmsg.GetErrMsg(code),
	})
}

// token 生成函数
func setToken(c *gin.Context, user models.User) {
	j := middleware.NewJWT()
	claims := middleware.MyClaims{
		Username: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.Error,
			"message": errmsg.GetErrMsg(errmsg.Error),
			"token":   token,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.Success,
		"id":      user.ID,
		"data":    user.Name,
		"message": errmsg.Success,
		"token":   token,
	})
	return
}
