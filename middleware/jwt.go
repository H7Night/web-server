package middleware

import (
	"errors"
	"net/http"
	"strings"
	"web-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		return errors.New("that's not even a token")
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return errors.New("token is expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return errors.New("couldn't handle this token")
	} else {
		return errors.New("token not active yet")
	}
}

// JwtToken jwt 中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		j := NewJWT()
		// 解析 token
		err := j.ParseToken(tokenString)
		if err != nil {
			var response gin.H
			if err == jwt.ErrSignatureInvalid {
				response = gin.H{"error": "Invalid token signature"}
			} else {
				response = gin.H{"error": "Failed to parse token"}
			}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		c.Next()
	}
}
