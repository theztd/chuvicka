package auth

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

func JWTValidate(tokenString string) bool {
	claimData := jwt.RegisteredClaims{}
	ct, err := jwt.ParseWithClaims(
		tokenString,
		&claimData,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTHash), nil
		},
	)
	if err != nil {
		log.Println("ERR: Token verificaition error", err)
		return false
	}

	fmt.Println("DEBUG: Parsed token (podpis): ", ct.Signature)
	fmt.Println("DEBUG: Parsed claim: ", claimData)

	if ct.Valid {
		return true
	}

	return false
}

func (usr *User) JWTGenerate() error {
	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chuvicka",
		ID:     usr.Email,
	})
	tokenStr, _ := _token.SignedString([]byte(JWTHash))
	usr.Token = tokenStr
	return nil
}
