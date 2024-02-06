package utils

import (
  "os"
  "time"
  "log"

  "github.com/golang-jwt/jwt"
)

func CreateToken(data string) (string, error) {
  token := jwt.NewWithClaims(
    jwt.SigningMethodHS256,
    jwt.MapClaims{
      "id": data,
      "exp": time.Now().Add(time.Hour * 24).Unix(),
    })

  return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseToken(token string) jwt.MapClaims {
  parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
  })

  if err != nil {
    log.Println(err)
  }

  tokenData := parsedToken.Claims.(jwt.MapClaims)

  return tokenData
}
