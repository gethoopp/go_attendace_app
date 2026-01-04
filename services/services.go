package services

import (
	"context"
	"database/sql"

	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/gethoopp/hr_attendance_app/database"
	"github.com/gethoopp/hr_attendance_app/modules"
	_ "github.com/go-sql-driver/mysql"
)

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
	var user modules.Users
	ctx := context.Background()

	userID, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer db.Close()

	query := `
		SELECT 
			id_users, 
			rfid_id, 
			id_first_name, 
			id_last_name, 
			id_departement, 
			email_user
		FROM Users 
		WHERE id_users = ?
	`

	row := db.QueryRowContext(ctx, query, userID)

	if err := row.Scan(
		&user.Id,
		&user.Rfid,
		&user.FirstName,
		&user.LastName,
		&user.Department,
		&user.Email,
	); err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "User tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal membaca data",
			"error":   err.Error(),
		})
		return
	}

	userResp := modules.UserResponse{
		Id:         user.Id,
		Rfid:       user.Rfid,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Department: user.Department,
		Email:      user.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Berhasil mengambil data user",
		"result":  userResp,
	})
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
