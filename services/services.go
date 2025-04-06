package services

import (
	"context"
	"database/sql"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"

	imageRecognition "github.com/gethoopp/hr_attendance_app/image_recognition"
	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gethoopp/hr_attendance_app/push_notification"
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server "})
		return
	}

	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	query := "INSERT INTO users (rfid_tag,name_user,email_user,divisi_user) VALUES(?,?,?,?)"

	_, err = db.ExecContext(ctx, query, log.Rfid, log.Name, log.Email, log.Divisi)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Sudah terdaftar"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Succes input data"})
	}

	msg_publish := fmt.Sprintf("Selamat Datang %s", log.Name)

	push_notification.Publisher_mssg(c, msg_publish)

}

func loginData(c *gin.Context) {
	var log modules.Userdata
	ctx := context.Background()
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "user login"})
	}

	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	query := "SELECT email_user, password_user FROM user_data "

	rows, err := db.QueryContext(ctx, query, log.Email)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"Status": http.StatusBadGateway, "Message": "Gagal melakukan request"})
	}

	defer rows.Close()

	for rows.Next() {
	}

}

func Register_Data(c *gin.Context) {
	ctx := context.Background()
	var log modules.RegUserdata

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Gagal menghubungkan ke database",
		})
	}

	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	resHash, err := bcrypt.GenerateFromPassword([]byte(log.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	query := "INSERT INTO regist_user(rfid_id,id_first_name,id_last_name,id_departement,email_user,password_user) VALUES (?,?,?,?,?,?)"

	_, err = db.ExecContext(ctx, query, log.Rfid, log.Name, log.LastName, log.Divisi, log.Email, resHash)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"Status": http.StatusBadGateway, "Message": "Gagal mengirimkan data"})
		return
	}

	//update table user
	queryS := "INSERT INTO users(rfid_id,id_first_name,id_departement,email_user)VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE rfid_id = VALUES(rid_id), id_first_name = VALUES(id_first_name), id_departement = VALUES(id_departement), email_user = VALUES(email_user)"

	_, err = db.ExecContext(ctx, queryS, log.Rfid, log.Name, log.LastName, log.Divisi, log.Email, resHash)
	if err != nil {

		c.JSON(http.StatusBadGateway, gin.H{"Status": http.StatusBadGateway, "Message": "Gagal Update data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": " berhasil dikirim"})
}

// func Attendace_user(c *gin.Context) {
// }

func CompareImageFromDb(c *gin.Context) {
	var log modules.DataImage
	ctx := context.Background()
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
		return
	}

	defer db.Close()
	query := "SELECT FROM * images_db where "
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal Server errror", "Status": http.StatusBadRequest})
		return
	}

	defer rows.Close()

	for rows.Next() {
		//

		//melakukan compare image
		imageRecognition.CompareImage(c, []byte(log.ImageRef), []byte(log.ImageRef))

	}

}
