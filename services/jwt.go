package services

import (
	"time"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

func CreateToken() (string, error) {

	var log modules.Users

	expiredTime := time.Now().Add(10 * time.Minute)

	claims := &modules.ClaimsData{
		NamaUser: log.FirstName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
