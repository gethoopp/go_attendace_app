package services

import (
	"net/http"
	"os"
	"time"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(user modules.Users) (string, error) {
	expiredTime := time.Now().Add(5 * time.Minute)

	claims := &modules.ClaimsData{
		UserID:   user.Id,
		NamaUser: user.FirstName,
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

//validate token is expired or not

func ValidateToken(token string, c *gin.Context) (*modules.ClaimsData, error) {
	claims := &modules.ClaimsData{}

	parsedToken, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"Status":  http.StatusBadRequest,
				"message": "Token Is Expired",
			},
		)
	}

	if !parsedToken.Valid {
		c.JSON(
			http.StatusConflict,
			gin.H{
				"Status":  http.StatusConflict,
				"message": "Token is not valid",
			},
		)
	}

	return claims, nil

}
