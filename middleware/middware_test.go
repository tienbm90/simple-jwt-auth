package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"testing"
)

func TestTokenVerify(t *testing.T) {
	//tokenString :="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI0NDQ0NDQ0IiwiZXhwIjoxNjA0OTg5MTgzLCJzdWIiOiJ0aWVuIn0.GxJrlMwOgrtL2fooE2D0K6qkeDv47QwX6BwBgy98cuzQ5lFiPmGgqOc7X05FIRzDKvljDCq0rpM19vzkVaypJQ"
	tokenString :="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI0NDQ0NDQ0IiwiZXhwIjoxNjA1MDAyMDY5LCJzdWIiOiJ0aWVuIn0.xp_dUN5LJgaTxCHHhX__hjEERV3xd1G5VjlGQRB1qI2yGidbevIXWXbY_QD0Jacl-BGdXOmoU3Jh_cpM4b2T3w"


	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("777"), nil
	})

	if err != nil {
		log.Println("wtf")
		log.Println(token)
		log.Println(err.Error())
	}else {
		print(token.SigningString())
	}
}