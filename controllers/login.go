package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"web-server/middleware"
	"web-server/models"
	"web-server/utils/errmsg"
)

func LoginFront(c *gin.Context) {
	var formData models.User
	_ = c.ShouldBind(&formData)
	var code int

	formData, code = models.CheckLoginFront(formData.Username, formData.Password)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    formData.Username,
		"id":      formData.ID,
		"message": errmsg.GetErrMsg(code),
	})
}

func Login(c *gin.Context) {
	var formData models.User
	_ = c.ShouldBind(&formData)
	var token string
	var code int

	formData, code = models.CheckLogin(formData.Username, formData.Password)

	if code == errmsg.SUCCESS {
		setToken(c, formData)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    formData.Username,
			"id":      formData.ID,
			"message": errmsg.GetErrMsg(code),
			"token":   token,
		})
	}
}

// token 生成函数
func setToken(c *gin.Context, user models.User) {
	j := middleware.NewJWT()
	claims := middleware.Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    "HE",
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"token":   token,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    user.Username,
		"id":      user.ID,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
		"token":   token,
	})
	return
}
