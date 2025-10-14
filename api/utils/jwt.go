package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(mail string, adm bool) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("server misconfigured: missing JWT_SECRET")
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mail": mail,
		"adm":  adm,
		"iat":  now.Unix(),
		"nbf":  now.Unix(),
		"exp":  now.Add(24 * time.Hour).Unix(),
		"iss":  "go-mongo-api",
	})
	return token.SignedString([]byte(secret))
}

func JWTValidate(tokenString string) (string, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return "", false, err
	}

	if !token.Valid {
		return "", false, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, errors.New("could not parse claims")
	}

	mail, ok := claims["mail"].(string)
	if !ok || mail == "" {
		return "", false, errors.New("mail not found in token")
	}

	adm, ok := claims["adm"].(bool)
	if !ok {
		return "", false, errors.New("adm status not found in token")
	}

	return mail, adm, nil
}
