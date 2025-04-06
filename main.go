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
	register := services.Register_Data

	//route URL
	r.POST("api/input", req)
	r.POST("api/data", func(ctx *gin.Context) {
		reqs(ctx)
	})
	r.POST("api/upload", res_image)
	r.GET("api/getData", compare_img)
	r.POST("api/register", register)

	// go func() {
	// 	fmt.Println("Profiling server running on :6060")
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
	r.Run(":8080")
}

//curl -X POST http://localhost:8080/api/register -d '{"rfid_id" : 12021,"id_first_name" : "Sidarrta", "id_last_name" : "andi", "id_departement" : "IT","email_user" : "andi@gmail.com", "password_user": "12322"}' -H "Content-Type: application/json"
