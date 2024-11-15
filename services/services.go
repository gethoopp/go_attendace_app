package services

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func User_data(c *gin.Context) {

	ctx := context.Background()

	if err := c.ShouldBindJSON(""); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	query := "INSERT INTO (rfid_tag) VALUES (?)"

	_, err = db.QueryContext(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
	}

}

func Attendace_user(c *gin.Context) {

}
