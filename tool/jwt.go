package tool

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWT struct {
	secret string // 密钥
	expire int64  // 过期时间间隔,单位秒
}

// NewJWT 创建jwt实例
func NewJWT(secret string, expire int64) *JWT {
	return &JWT{
		secret: secret,
		expire: expire,
	}
}

// GenerateToken 生成token
func (j *JWT) GenerateToken(claims map[string]interface{}) (string, error) {
	claims["exp"] = time.Now().Add(time.Second * time.Duration(j.expire)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString([]byte(j.secret))
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
