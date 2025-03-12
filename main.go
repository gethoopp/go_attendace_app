package main

import (
	imageRecognition "github.com/gethoopp/hr_attendance_app/image_recognition"
	"github.com/gethoopp/hr_attendance_app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	req := services.Input_rfid
	reqs := services.User_data
	res_image := imageRecognition.Save_image

	r.POST("api/input", req)
	r.POST("api/data", reqs)
	r.POST("api/upload", res_image)

	r.Run(":8080")
	imageRecognition.OpenVideo()
}

//curl -X POST http://localhost:8080/api/data -d '{"nama_user" : "andi" : "email_user" : andi@gmail.com : "divisi_user" : "IT"}' -H "Content-Type: application/json"
