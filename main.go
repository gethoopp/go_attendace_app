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
	presence := services.Get_presence
	getPresenceBydate := services.Get_presence_byDate
	register := services.Register_Data
	login := services.LoginData
	logout := services.Logout_User
	saveToken := push_notification.SaveDeviceToken
	chatUser := chat.ChatBotOllama
	JWTMiddleware := middleware.JWTMiddleware()
	sendNotif := push_notification.SendsNotification
	checkIn := services.Check_in
	checkOut := services.Check_out
	initFirebase := middleware.InitFirebase

	// ===== ROUTES =====
	r.GET("/ws/input", req)
	r.GET("/api/data", JWTMiddleware, reqs)
	r.GET("/api/presence", JWTMiddleware, presence)
	r.POST("/api/getByDate", JWTMiddleware, getPresenceBydate)
	r.POST("/api/register", register)
	r.POST("/api/login", login)
	r.POST("/api/logout", JWTMiddleware, logout)
	r.POST("/api/saveToken", JWTMiddleware, saveToken)
	r.POST("/api/sendNotif", initFirebase, sendNotif)
	r.POST("/api/chat", JWTMiddleware, chatUser)
	r.POST("/api/checkIn", JWTMiddleware, checkIn)
	r.PUT("/api/checkOut", JWTMiddleware, checkOut)

	// ===== PORT HEROKU =====
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
