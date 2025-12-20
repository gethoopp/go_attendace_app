package main

import (
	"os"

	"github.com/gethoopp/hr_attendance_app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Pakai release mode (lebih hemat RAM)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ===== ROUTE WAJIB UNTUK HEROKU =====
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    "hr-attendance",
		})
	})

	// ===== SERVICE =====
	req := services.Input_rfid
	reqs := services.User_data
	register := services.Register_Data
	login := services.LoginData

	// ===== ROUTES =====
	r.GET("/ws/input", req)
	r.POST("/api/data", reqs)
	r.POST("/api/register", register)
	r.POST("/api/login", login)

	// ===== PORT HEROKU =====
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
