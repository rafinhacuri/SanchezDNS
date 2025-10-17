package passwords

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func BCrypt(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return "{BCRYPT}" + string(hashed), nil
}

func VerifyBCrypt(password, hashed string) bool {
	hashed = strings.TrimPrefix(hashed, "{BCRYPT}")
	if strings.HasPrefix(hashed, "{CRYPT}") {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
