package services

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gethoopp/hr_attendance_app/database"
	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// function logout
func Logout_User(c *gin.Context) {
	var userResp modules.Users
	ctx := context.Background()

	db, err := database.GetDB()
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

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghubungkan ke database",
			"error":   err.Error(),
		})
		return
	}
	defer db.Close()

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

	db, err := database.GetDB()
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

		switch {
		case user.Email == res.Email:
			c.JSON(http.StatusConflict, gin.H{
				"message": "Email sudah terdaftar",
			})
			return
		case user.Rfid == res.Rfid:
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
