package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"goIland/utils"
)

type JWT struct {
	JwtKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(utils.JwtKey),
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	TokenExpired     = errors.New("1")
	TokenNotValidYet = errors.New("2")
	TokenMalformed   = errors.New("3")
	TokenInvalid     = errors.New("4")
)

func (j *JWT) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}
