package tool

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	j := NewJWT("secret", 5)
	token, err := j.GenerateToken(map[string]interface{}{
		"username": "admin",
		"id":       1,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
	time.Sleep(time.Second * 10)
	claims, err := j.ParseToken(token)
	if err != nil {
		t.Error(err)
	}

	t.Log(claims)
}
