package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"web-server/utils"
	"web-server/utils/errmsg"
)

type JWT struct {
	JWTKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		JWTKey: []byte(utils.JwtKey),
	}
}

type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

// CreateToken sha256生成Token
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JWTKey)
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.JWTKey, nil
	})
	// 验证token
	if token.Valid {
		return nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return TokenMalformed
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return TokenExpired
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return TokenInvalid
	} else {
		return TokenNotValidYet
	}
}

// JwtToken jwt 中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			code = errmsg.ErrorTokenExist
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		checkToken := strings.Split(tokenString, " ")
		if len(checkToken) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		j := NewJWT()
		// 解析 token
		err := j.ParseToken(checkToken[1])
		if err != nil {
			if errors.Is(err, TokenExpired) {
				c.JSON(http.StatusOK, gin.H{
					"status":  errmsg.Error,
					"message": "token is expired",
					"data":    nil,
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.Error,
				"message": err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
