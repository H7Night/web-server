package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"web-server/utils"
	"web-server/utils/errmsg"
)

type JWT struct {
	JwtKey []byte
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWT() *JWT {
	return &JWT{
		JwtKey: []byte(utils.JwtKey),
	}
}

func (j *JWT) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ErrorTokenExist,
				"message": errmsg.GetErrMsg(errmsg.ErrorTokenExist),
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(tokenHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ErrorTokenTypeWrong,
				"message": errmsg.GetErrMsg(errmsg.ErrorTokenTypeWrong),
			})
			c.Abort()
			return
		}
		claims, err := NewJWT().ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ErrorTokenWrong,
				"message": errmsg.GetErrMsg(errmsg.ErrorTokenWrong),
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ErrorTokenRuntime,
				"message": errmsg.GetErrMsg(errmsg.ErrorTokenRuntime),
			})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}
