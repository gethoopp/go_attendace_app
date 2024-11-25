package main

import (
	"github.com/gethoopp/hr_attendance_app/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	req := services.User_data

	r.POST("api/input", req)

	r.Run(":8080")
}

//curl -X POST http://localhost:8080/api/input -d '{"rfid_tag" : 1234589}' -H "Content-Type: application/json"
