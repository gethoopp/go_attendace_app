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

	query := "SELECT id_users, rfid_id ,id_first_name, id_last_name,  id_departement, email_user,password_user FROM Users WHERE id_users=?"

	rows, err := db.QueryContext(ctx, query, log.Rfid, log.FirstName, log.LastName, log.Department, log.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Sudah terdaftar"})
	}

	if !rows.Next() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "ID tidak ditemukan",
		})
		return
	}

	if err := rows.Scan(&log.Id, &log.Rfid, &log.FirstName, &log.LastName, &log.Department, &log.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membaca data",
			"error":   err.Error(),
		})
		return
	}

	userResp := modules.UserResponse{
		Id:         log.Id,
		Rfid:       log.Rfid,
		FirstName:  log.FirstName,
		LastName:   log.LastName,
		Department: log.Department,
		Email:      log.Email,
	}

	// msg_publish := fmt.Sprintf("Selamat Datang %s", log.FirstName)
	// push_notification.Publisher_mssg(c, msg_publish)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Login berhasil",
		"result":  userResp,
	})

}

func LoginData(c *gin.Context) {
	var log modules.Users
	var reqLogin modules.LoginRequest
	ctx := context.Background()

	// Tampilkan error detail agar mudah debug
	if err := c.ShouldBindJSON(&reqLogin); err != nil {
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

	// var emailUser string
	// var passwordHash string
	// var idUser int

	query := "SELECT id_users, rfid_id ,id_first_name, id_last_name,  id_departement, email_user,password_user FROM Users WHERE email_user=?"

	rows, err := db.QueryContext(ctx, query, reqLogin.Email)
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
	if err := rows.Scan(&log.Id, &log.Rfid, &log.FirstName, &log.LastName, &log.Department, &log.Email, &log.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membaca data",
			"error":   err.Error(),
		})
		return
	}

	// Cek password hash
	if err := bcrypt.CompareHashAndPassword([]byte(log.Password), []byte(reqLogin.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Password salah",
		})
		return
	}

	// Buat token
	token, err := CreateToken(
		log,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
		return
	}

	userResp := modules.UserResponse{
		Id:         log.Id,
		Rfid:       log.Rfid,
		FirstName:  log.FirstName,
		LastName:   log.LastName,
		Department: log.Department,
		Email:      log.Email,
	}

	// Sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Login berhasil",
		"token":   token,
		"result":  userResp,
	})
}

func Register_Data(c *gin.Context) {
	ctx := context.Background()

	var user modules.Users
	var res modules.UserResponse

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal koneksi database",
		})
		return
	}
	defer db.Close()

	queryCheck := `
		SELECT rfid_id, email_user
		FROM Users
		WHERE email_user = ? OR rfid_id = ?
	`

	err = db.QueryRowContext(ctx, queryCheck, user.Email, user.Rfid).
		Scan(&res.Rfid, &res.Email)

	if err == nil {

		if user.Email == res.Email {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Email sudah terdaftar",
			})
			return
		}
		if user.Rfid == res.Rfid {
			c.JSON(http.StatusConflict, gin.H{
				"message": "RFID sudah terdaftar",
			})
			return
		}
	}

	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal cek data",
			"error":   err.Error(),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal hash password",
		})
		return
	}

	queryInsert := `
		INSERT INTO Users
		(rfid_id, id_first_name, id_last_name, id_departement, email_user, password_user)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = db.ExecContext(
		ctx,
		queryInsert,
		user.Rfid,
		user.FirstName,
		user.LastName,
		user.Department,
		user.Email,
		hash,
	)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Gagal menyimpan data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
	})
}

func Logout_User(c *gin.Context) {
	var userResp modules.Users
	ctx := context.Background()

	db, err := GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal koneksi database",
		})
		return
	}

	defer db.Close()

	query := "UPDATE Users SET deleted_add = NOW() WHERE rfid_id = ? AND deleted_add IS NULL"

	rows, err := db.ExecContext(ctx, query, userResp.Rfid)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Gagal logout user",
			"error":   err.Error(),
		})
		return
	}

	result, _ := rows.RowsAffected()
	if result == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User tidak ditemukan atau sudah logout",
		})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "berhasil logout",
		},
	)
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
