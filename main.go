package main

import (
	"github.com/gethoopp/hr_attendance_app/chat"
	"github.com/gethoopp/hr_attendance_app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// gin.SetMode(gin.ReleaseMode)

	req := services.Input_rfid
	reqs := services.User_data
	// res_image := imageRecognition.Save_image
	// compare_img := services.CompareImageFromDb
	register := services.Register_Data
	chat_regist := chat.ChatBotOllama
	login := services.LoginData

	//route URL
	r.POST("api/input", req)
	r.POST("api/data", func(ctx *gin.Context) {
		reqs(ctx)
	})
	// r.POST("api/upload", res_image)
	// r.GET("api/getData", compare_img)
	r.POST("api/register", register)
	r.POST("api/chat/chatbot", chat_regist)
	r.POST("api/login", login)

	r.Run(":8080")

}

//curl -X POST http://localhost:8080/api/register -d '{"rfid_id" : 12021,"id_first_name" : "Sidarrta", "id_last_name" : "andi", "id_departement" : "IT","email_user" : "andi@gmail.com", "password_user": "12322"}' -H "Content-Type: application/json"
//Test chat bot
// curl -X POST http://localhost:8080/api/chat/chatbot \
//   -H "Content-Type: application/json" \
//   -d '{"prompt": "Saya mencintai seseorang yang sudah saya kenal setahunan ini, bagaiaman cara menghilangkpan perasaan itu?"}'

// curl: (7) Failed to connect to localhost port 8080 after 0 ms: Couldn't connect to server
//
