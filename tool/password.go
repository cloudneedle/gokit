package tool

import "golang.org/x/crypto/bcrypt"

type Password string

func (p *Password) String() string {
	return string(*p)
}

// Encrypt 加密密码
func (p *Password) Encrypt() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.String()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*p = Password(hash)
	return nil
}

// Check 校验密码
func (p *Password) Check(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.String()), []byte(password))
	if err != nil {
		return false
	}
	return true
}
