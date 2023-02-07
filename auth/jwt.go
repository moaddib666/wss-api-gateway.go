package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
)

var appSecret = []byte(os.Getenv("APP_SECRET"))

type JWTClient struct {
	token  *jwt.Token
	claims jwt.MapClaims
}

func (c *JWTClient) Id() string {
	return c.claims["jti"].(string)
}

func (c *JWTClient) IsAuthorised() bool {
	return c.token.Valid
}

func (c *JWTClient) IsAdmin() bool {
	return false
}

func (c *JWTClient) init(header []byte) {
	var err error
	if len(header) < 8 {
		log.Printf("invalid auth header `%s` reciaved, connection dropped", header)
		return
	}
	tokenText := string(header[7:]) // Remove Bearer
	c.token, err = jwt.Parse(tokenText, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return appSecret, nil
	})
	c.claims, _ = c.token.Claims.(jwt.MapClaims)
	if err != nil {
		log.Printf("Can't verify security token `%+v`, connection dropped", err)
		return
	}
}
