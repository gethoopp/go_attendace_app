package push_notification

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"firebase.google.com/go/v4/messaging"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gin-gonic/gin"
)

func GetDB() (*sql.DB, error) {
	// 1. PRIORITY: Heroku (JAWSDB)
	if jaws := os.Getenv("JAWSDB_URL"); jaws != "" {
		u, err := url.Parse(jaws)
		if err != nil {
			return nil, err
		}

		user := u.User.Username()
		pass, _ := u.User.Password()
		host := u.Host
		dbname := strings.TrimPrefix(u.Path, "/")

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}

		db.SetMaxIdleConns(20)
		db.SetMaxOpenConns(10)
		db.SetConnMaxLifetime(time.Minute * 5)

		return db, nil
	}

	// 2. fallback ke env lokal
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	return db, nil
}

func SaveDeviceToken(c *gin.Context) {
	var deviceToken modules.DeviceTokenRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&deviceToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//buka connection
	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghubungkan ke database",
			"error":   err.Error(),
		})
		return
	}

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	userClaims, ok := claims.(*modules.ClaimsData)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token claims",
		})
		return
	}
	UserId := userClaims.UserID

	defer db.Close()

	query := "INSERT INTO device_tokens(user_id,device_token,platform) VALUES (?,?,?)  ON DUPLICATE KEY UPDATE device_token = ?"

	_, err = db.ExecContext(ctx, query, UserId, deviceToken.DeviceToken, deviceToken.Platform, deviceToken.DeviceToken)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Gagal melakukan request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Berhasil Mendapatkan token",
		"status":  http.StatusOK,
		"result":  userClaims.ExpiresAt,
	})

}

func SendsNotification(c *gin.Context) {
	ctx := context.Background()
	fcmAny, ok := c.Get("fcm")
	if !ok {
		c.JSON(500, gin.H{"message": "FCM client not found"})
		return
	}
	var request modules.Message

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	message := &messaging.Message{
		Token: request.Token,
		Notification: &messaging.Notification{
			Title: request.Title,
			Body:  request.Body,
		},
	}

	_, err := fcmAny.(*messaging.Client).Send(ctx, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Firebase gagal",
		})
		return
	}

}
