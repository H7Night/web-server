package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"web-server/middleware"
	"web-server/models"
	"web-server/utils/errmsg"
)

func Login(c *gin.Context) {
	var formData models.User
	_ = c.ShouldBind(&formData)
	var code int

	formData, code = models.CheckLogin(formData.Name, formData.Password, c.Param("state"))
	if code == errmsg.Success {
		setToken(c, formData)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    formData.Name,
		"id":      formData.ID,
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
