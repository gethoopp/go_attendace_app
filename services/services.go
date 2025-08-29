package services

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gethoopp/hr_attendance_app/push_notification"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
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
	db.SetConnMaxLifetime(10)

	return db, nil
}

func Input_rfid(c *gin.Context) {

	ctx := context.Background()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
	}

	defer db.Close()

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
	var log modules.Users
	ctx := context.Background()
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}

	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server "})
		return
	}

	defer db.Close()

	query := "INSERT INTO users (rfid_tag,name_user,email_user,divisi_user) VALUES(?,?,?,?)"

	_, err = db.ExecContext(ctx, query, log.Rfid, log.FirstName, log.Email, log.Department)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Sudah terdaftar"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Succes input data"})
	}

	msg_publish := fmt.Sprintf("Selamat Datang %s", log.FirstName)

	push_notification.Publisher_mssg(c, msg_publish)

}

func LoginData(c *gin.Context) {
	var log modules.Users
	ctx := context.Background()
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}
	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "user login"})
	}

	defer db.Close()

	var email_user string
	var password_user string

	query := "SELECT email_user, password_user FROM user_data WHERE email_user = ?"

	rows, err := db.QueryContext(ctx, query, log.Email)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"Status": http.StatusBadGateway, "Message": "Gagal melakukan request"})
	}

	defer rows.Close()
	err = rows.Scan(&email_user, &password_user)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Terjadi kesalahan internal"})
		return
	}

	//compare hash password
	err = bcrypt.CompareHashAndPassword([]byte(password_user), []byte(log.Password))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Masukkan informasi login dnegan benar",
		})
		return
	}

}

func Register_Data(c *gin.Context) {
	ctx := context.Background()
	var log modules.Users

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}

	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Gagal menghubungkan ke database",
		})
	}

	defer db.Close()

	resHash, err := bcrypt.GenerateFromPassword([]byte(log.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	query := "INSERT INTO Users(rfid_id,id_first_name,id_last_name,id_departement,email_user,password_user) VALUES (?,?,?,?,?,?)"

	_, err = db.ExecContext(ctx, query, log.Rfid, log.FirstName, log.LastName, log.Department, log.Email, resHash)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Status":  http.StatusBadGateway,
			"Message": "Gagal mengirimkan data",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " berhasil dikirim"})
}

// func Attendace_user(c *gin.Context) {
// }

// func CompareImageFromDb(c *gin.Context) {
// 	var log modules.DataImage
// 	ctx := context.Background()
// 	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server errror"})
// 		return
// 	}

// 	defer db.Close()
// 	query := "SELECT FROM * images_db where "
// 	rows, err := db.QueryContext(ctx, query)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal Server errror", "Status": http.StatusBadRequest})
// 		return
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		//

// 		//melakukan compare image
// 		// imageRecognition.CompareImage(c, []byte(log.ImageRef), []byte(log.ImageRef))

// 	}

// }
