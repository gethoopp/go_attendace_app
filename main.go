package main

import (
	"os"

	"github.com/gethoopp/hr_attendance_app/chat"
	"github.com/gethoopp/hr_attendance_app/middleware"
	"github.com/gethoopp/hr_attendance_app/push_notification"
	"github.com/gethoopp/hr_attendance_app/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Pakai release mode (lebih hemat RAM)
	initFirebase := middleware.InitFirebase
	godotenv.Load()
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
	logout := services.Logout_User
	saveToken := push_notification.SaveDeviceToken
	chatUser := chat.ChatBotOllama
	JWTMiddleware := middleware.JWTMiddleware()
	sendNotif := push_notification.SendsNotification

	// ===== ROUTES =====
	r.GET("/ws/input", req)
	r.POST("/api/data", JWTMiddleware, reqs)
	r.POST("/api/register", register)
	r.POST("/api/login", login)
	r.POST("/api/logout", JWTMiddleware, logout)
	r.POST("/api/saveToken", JWTMiddleware, saveToken)
	r.POST("/api/sendNotif", sendNotif, initFirebase)
	r.POST("/api/chat", JWTMiddleware, chatUser)

	// ===== PORT HEROKU =====
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
