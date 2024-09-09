package api

import (
	"fmt"
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
	var state string
	//登录校验
	formData, state, code = models.CheckLogin(formData.Name, formData.Password)
	if code == errmsg.Success {
		setToken(c, formData, state)
	} else {
		fmt.Println("login faild")
	}
}

// token 生成函数
func setToken(c *gin.Context, user models.User, state string) {
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
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.Success,
		"id":      user.ID,
		"state":   state,
		"data":    user.Name,
		"message": errmsg.Success,
		"token":   token,
	})
	return
}
