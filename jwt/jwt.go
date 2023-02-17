package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken 生成token
func GenerateToken(key string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

// ParseToken 解析 token
func ParseToken(key, tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return err
	}

	if token.Valid {
		return nil
	} else {
		return err
	}
}
