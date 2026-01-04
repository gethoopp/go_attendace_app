package services

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gethoopp/hr_attendance_app/database"
	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gin-gonic/gin"
)

// function check in
func Check_in(c *gin.Context) {
	var log modules.Attendance
	ctx := context.Background()

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  http.StatusBadGateway,
				"message": "Gagal membaca data",
			},
		)
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"status":  http.StatusForbidden,
				"message": "Gagal menghubungkan database",
			},
		)
		return
	}

	defer db.Close()

	query := "SELECT * FROM attendances WHERE user_id = ? AND attendance_date = CURDATE() AND check_out IS NULL"

	err = db.QueryRowContext(ctx, query, log.IdUser).Scan(&log.Id, &log.IdUser, &log.Check_in, &log.Check_out, &log.DateIn, &log.Status, &log.Created_at)

	if err == sql.ErrNoRows {
		query := "INSERT INTO attendances (user_id, check_in, attendance_date) VALUES (?, NOW(), CURDATE())"

		_, err = db.ExecContext(
			ctx, query, log.IdUser,
		)
		if err != nil {
			c.JSON(
				http.StatusFound,
				gin.H{
					"status":  http.StatusFound,
					"message": "Gagal Check in",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Check in berhasil",
			},
		)

		return

	}

}

func Check_out(c *gin.Context) {
	var log modules.Attendance
	ctx := context.Background()

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  http.StatusBadGateway,
				"message": "",
			},
		)
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Gagal menghubungkan database",
			},
		)

		return
	}

	defer db.Close()

	query := "UPDATE attendances SET check_out = NOW(), status = 'OUT' WHERE user_id = ? AND attendance_date = CURDATE() AND check_out IS NULL;"
	_, err = db.ExecContext(ctx, query, log.IdUser)
	if err != nil {
		c.JSON(
			http.StatusFound,
			gin.H{
				"status":  http.StatusFound,
				"message": "Data tidak ditemukan",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Check out berhasil",
		},
	)

}
