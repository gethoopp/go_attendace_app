package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/gethoopp/hr_attendance_app/modules"
	_ "github.com/go-sql-driver/mysql"
)

func Input_rfid(c *gin.Context) {

	ctx := context.Background()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
	}

	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	conn, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ada Kesalahaan Internal"})
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "WebSocket handshake failed",
				"detail": err.Error(),
			})
			return
		}

		var value int
		_, err = fmt.Sscanf(string(message), "%d", &value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "WebSocket failed request",
				"detail": err.Error(),
			})
			return
		}

		query := "INSERT INTO users(rfid_tag) VALUES(?)"

		_, err = db.ExecContext(ctx, query, value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Succes input data"})
		}

	}

}

func User_data(c *gin.Context) {
	var log modules.Userdata
	ctx := context.Background()
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
		return
	}

	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	//insert data from flutter input
	query := "INSERT INTO users (name_user,email_user,divisi_user) VALUES(?,?,?)"

	_, err = db.ExecContext(ctx, query, log.Name, log.Email, log.Divisi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Succes input data"})
	}

}

// func Attendace_user(c *gin.Context) {

// }

// func CompareImageFromDb(c *gin.Context) {

// }
