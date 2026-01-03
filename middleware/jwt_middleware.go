package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gethoopp/hr_attendance_app/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header kosong",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Format Authorization salah",
			})
			return
		}

		tokenStr := parts[1]
		secretKey := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.ParseWithClaims(
			tokenStr,
			&modules.ClaimsData{},
			func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return secretKey, nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak valid",
			})
			return
		}

		claims, ok := token.Claims.(*modules.ClaimsData)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Claims tidak valid",
			})
			return
		}

		_, err = services.ValidateToken(
			tokenStr, &gin.Context{},
		)

		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Token tidak ditemukan",
				},
			)
		}

		c.Set("claims", claims)
		c.Set("id_user", claims.UserID)

		c.Next()
	}
}
