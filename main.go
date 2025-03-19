package main

import (
	imageRecognition "github.com/gethoopp/hr_attendance_app/image_recognition"
	"github.com/gethoopp/hr_attendance_app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// gin.SetMode(gin.ReleaseMode)

	req := services.Input_rfid
	reqs := services.User_data
	res_image := imageRecognition.Save_image
	compare_img := services.CompareImageFromDb

	r.POST("api/input", req)
	r.POST("api/data", func(ctx *gin.Context) {
		reqs(ctx)
	})
	r.POST("api/upload", res_image)
	r.GET("api/getData", compare_img)

	r.Run(":8080")
	// imageRecognition.OpenVideo()
}

//curl -X POST http://localhost:8080/api/data -d '{"rfid_tag" : 12021,"nama_user" : "Sidarta", "email_user" : "andi@gmail.com", "divisi_user" : "IT"}' -H "Content-Type: application/json"
