package services

import (
	"context"
	"database/sql"
	"net/http"
	"time"

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

	query := "SELECT * FROM attendances WHERE user_id = ? AND attendance_date = CURDATE() AND status = 'OUT' "

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

	queries := "UPDATE attendances SET check_in = NOW(), status = 'IN', check_out = NULL WHERE user_id = ? AND attendance_date = CURDATE()"
	_, err = db.ExecContext(ctx, queries, log.IdUser)
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
			"message": "User sudah check in",
		},
	)

}

func Get_presence_byDate(c *gin.Context) {
	var log modules.Attendance
	var req modules.AttendanceRequest

	ctx, cancel := context.WithTimeout(c, 100*time.Millisecond)
	defer cancel()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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

	t, err := time.Parse(time.RFC3339Nano, req.DateIn)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  http.StatusBadGateway,
				"message": "Gagal parsing data",
			},
		)
		return
	}

	queries := "SELECT id_attendance, user_id, check_in, check_out, attendance_date, status, created_at FROM attendances WHERE user_id = ? AND attendance_date = ?"
	rows, err := db.QueryContext(ctx, queries, req.IdUser, t)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "Gagal melakukan request",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Data tidak ditemukan",
		})
		return
	}

	if err := rows.Scan(
		&log.Id,
		&log.IdUser,
		&log.Check_in,
		&log.Check_out,
		&log.DateIn,
		&log.Status,
		&log.Created_at,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membaca data",
			"error":   err.Error(),
		})
		return

	}

	userResp := modules.Attendance{
		Id:         log.Id,
		IdUser:     log.IdUser,
		Check_in:   log.Check_in,
		Check_out:  log.Check_out,
		DateIn:     log.DateIn,
		Status:     log.Status,
		Created_at: log.Created_at,
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Berhasil mendapatkan data",
			"result":  userResp,
		},
	)

}

func Check_out(c *gin.Context) {
	var log modules.Attendance
	ctx := context.Background()

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Request body tidak valid",
		})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghubungkan database",
		})
		return
	}
	defer db.Close()

	query := `
		UPDATE attendances
		SET check_out = NOW(), status = 'OUT'
		WHERE user_id = ?
		  AND attendance_date = CURDATE()
		  AND check_out IS NULL
	`

	result, err := db.ExecContext(ctx, query, log.IdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal melakukan check out",
		})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User belum check-in atau sudah check-out",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Check out berhasil",
	})
}

func Get_presence(c *gin.Context) {
	var log modules.Attendance

	ctx := context.Background()

	userID, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter date wajib diisi (YYYY-MM-DD)",
		})
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

	t, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  http.StatusBadGateway,
				"message": "Gagal parsing data",
			},
		)
		return
	}

	queries := "SELECT id_attendance, user_id, check_in, check_out, attendance_date, status, created_at FROM attendances WHERE user_id = ? AND attendance_date = ? "

	rows := db.QueryRowContext(ctx, queries, userID, t)

	if err := rows.Scan(
		&log.Id,
		&log.IdUser,
		&log.Check_in,
		&log.Check_out,
		&log.DateIn,
		&log.Status,
		&log.Created_at); err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "User tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Gagal membaca data",
			"error":   err.Error(),
		})
		return
	}

	userResp := modules.Attendance{
		Id:         log.Id,
		IdUser:     log.IdUser,
		Check_in:   log.Check_in,
		Check_out:  log.Check_out,
		DateIn:     log.DateIn,
		Status:     log.Status,
		Created_at: log.Created_at,
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Berhasil mendapatkan data",
			"result":  userResp,
		},
	)

}
