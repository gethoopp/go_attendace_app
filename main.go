package main

import (
	"strings"

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

	// Test chatbot with spago
	// fmt.Println("Chatbot activated")
	// scanner := bufio.NewScanner(os.Stdin)

	// for {
	// 	fmt.Print("Kamu: ")
	// 	scanner.Scan()
	// 	userInput := scanner.Text()

	// 	if strings.ToLower(userInput) == "keluar" {
	// 		fmt.Println("Chatbot: Sampai jumpa! 👋")
	// 		break
	// 	}

	// 	response := chatbotResponse(userInput)
	// 	fmt.Println("Chatbot:", response)
	// }

}

func chatbotResponse(inputText string) string {
	switch {
	case strings.Contains(inputText, "halo"):
		return "Halo juga! Ada yang bisa saya bantu?"
	case strings.Contains(inputText, "siapa kamu"):
		return "Saya adalah chatbot buatan kamu sendiri 😄"
	case strings.Contains(inputText, "terima kasih"):
		return "Sama-sama! Senang bisa membantu 👍"
	case strings.Contains(inputText, "cuaca"):
		return "Maaf, saya belum bisa memberikan info cuaca saat ini."
	default:
		return "Maaf, saya belum mengerti maksudmu. Coba pertanyaan lain?"
	}
}

//curl -X POST http://localhost:8080/api/register -d '{"rfid_id" : 12021,"id_first_name" : "Sidarrta", "id_last_name" : "andi", "id_departement" : "IT","email_user" : "andi@gmail.com", "password_user": "12322"}' -H "Content-Type: application/json"
