package util

import (
	"github.com/dgrijalva/jwt-go"
)

/**
 * 解析 token
 */
func ParseToken(tokenSrt string, SecretKey []byte) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return
	}
	claims = token.Claims
	return
}
