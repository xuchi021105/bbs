package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SigningKey = []byte("xuchi") // 用于jwt加密的签名(或者说密钥)

type AuthClaims struct {
	jwt.RegisteredClaims
	Uid uint `json:"uid"`
}

const TokenExpireDuration = time.Hour * 24

func GeneToken(uid uint) (string, error) {

	authClaims := AuthClaims{
		Uid: uid,
	}

	authClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(TokenExpireDuration))

	authClaims.Issuer = "xuchi"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	tokenString, err := token.SignedString(SigningKey)

	return "Bearer " + tokenString, err
}

func ParseToken(tokenString string) (*AuthClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invaild token")

}
