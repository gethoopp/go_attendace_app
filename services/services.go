package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gethoopp/hr_attendance_app/push_notification"
	_ "github.com/go-sql-driver/mysql"
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

		db.SetMaxIdleConns(5)
		db.SetMaxOpenConns(5)
		db.SetConnMaxLifetime(time.Minute * 3)

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

func Input_rfid(c *gin.Context) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		var msg modules.Users

		if err := conn.ReadJSON(&msg.Rfid); err != nil {
			log.Println("ReadJSON error:", err)
			break
		}

		log.Println("RFID:", msg.Rfid)

		// echo balik
		if err := conn.WriteJSON(gin.H{
			"rfid_tag": msg.Rfid,
		}); err != nil {
			log.Println("WriteJSON error:", err)
			conn.WriteJSON(gin.H{
				"error":  "Internal Server Error",
				"Detail": err.Error(),
			})
			break
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

	// Tampilkan error detail agar mudah debug
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghubungkan ke database",
			"error":   err.Error(),
		})
		return
	}
	defer db.Close()

	var emailUser string
	var passwordHash string

	query := "SELECT email_user, password_user FROM Users WHERE email_user = ?"

	rows, err := db.QueryContext(ctx, query, log.Email)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Gagal melakukan request",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	// HARUS rows.Next() dulu
	if !rows.Next() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Email atau password salah",
		})
		return
	}

	// Baru Scan
	if err := rows.Scan(&emailUser, &passwordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membaca data",
			"error":   err.Error(),
		})
		return
	}

	// Cek password hash
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(log.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Password salah",
		})
		return
	}

	// Buat token
	token, err := CreateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
		return
	}

	// Sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Login berhasil",
		"token":   token,
	})
}

func Register_Data(c *gin.Context) {
	ctx := context.Background()
	var log modules.Users

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}, "status": err.Error()})
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
