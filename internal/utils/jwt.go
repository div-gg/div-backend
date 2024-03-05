package utils

import (
	"log"
	"time"

  "github.com/divinitymn/div-backend/internal/config"

	"github.com/golang-jwt/jwt"
)

func CreateToken(data string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  data,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

  return token.SignedString([]byte(config.Env.JWTSecret))
}

func ParseToken(token string) jwt.MapClaims {
	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Env.JWTSecret), nil
	})

	if err != nil {
		log.Println(err)
	}

	tokenData := parsedToken.Claims.(jwt.MapClaims)

	return tokenData
}
