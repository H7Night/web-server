package api

import (
	"fmt"
	"net/http"
	"time"
	"web-server/middleware"
	"web-server/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/pkg/errors"
)

// Login 登录，
func Login(c *gin.Context) {
	var formData models.User
	if err := c.ShouldBind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "invalid request data")})
		return
	}

	var user models.User
	var state string
	user, state, err := models.CheckLogin(formData.Name, formData.Password)
	if err == nil {
		setToken(c, user, state)
	} else {
		fmt.Println("Login failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
	}
}

// token 生成函数
func setToken(c *gin.Context, user models.User, state string) {
	jwtMiddleware := middleware.NewJWT()

	claims := middleware.MyClaims{
		Username: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwtMiddleware.CreateToken(claims)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"id":      user.ID,
		"state":   state,
		"data":    user.Name,
		"message": "success",
		"token":   token,
	})
}
