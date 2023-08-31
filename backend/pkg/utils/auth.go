package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(userid int) (string, string, error) {
	uuid := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userid,
		"uuid":   uuid,
		"iss":    "fakebook",
	})

	tkn, err := token.SignedString([]byte(GetEnv("SECRETKEY")))
	if err != nil {
		return "", "", fmt.Errorf("generating JWT: %v", err)
	}

	return tkn, uuid, nil
}

func GetClaimsFromJWT(tkn string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetEnv("SECRETKEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token or claims")
	}
}

func HashPassword(password string) (string, error) {
	hashedpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedpw), nil
}
