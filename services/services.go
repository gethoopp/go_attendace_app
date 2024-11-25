package services

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gethoopp/hr_attendance_app/modules"
	_ "github.com/go-sql-driver/mysql"
)

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
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	//insert rfid number from scan to database
	query := "INSERT INTO users(rfid_tag) VALUES (?)"

	_, err = db.ExecContext(ctx, query, log.Rfid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Succes input data"})
	}

}

func Attendace_user(c *gin.Context) {

}
