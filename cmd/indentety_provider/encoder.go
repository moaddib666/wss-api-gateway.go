package main

import (
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func main() {
	var user string
	flag.StringVar(&user, "user", "", "# Pass user identity like: '`JohnSnow`'")
	flag.Parse()

	if len(user) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var sampleSecretKey = []byte(os.Getenv("MARGAY_AUTH_SECRET"))
	if len(sampleSecretKey) == 0 {
		fmt.Println("Please set secret variable, example: `export MARGAY_AUTH_SECRET=SuperSecret`")
		os.Exit(1)
	}
	claims := jwt.StandardClaims{
		Audience:  "localhost",
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Id:        user,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "IdentityServiceLocal",
		Subject:   "client",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		panic(err)
	}

	decodedToken, _ := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return sampleSecretKey, nil
	})
	fmt.Printf("Decoded token: `%+v` \n", decodedToken)
	fmt.Println("---")
	fmt.Printf("Raw Token: '%s'\n", tokenString)
}
