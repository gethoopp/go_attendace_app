package imageRecognition

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

func OpenVideo() {
	deviceID := 0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

}

func Save_image(c *gin.Context) {
	imageResult, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Gagal Mengirimkan gambar ke server"})
		return
	}

	file, err := imageResult.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": http.StatusInternalServerError, "message": "Gagal membuka file"})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": http.StatusInternalServerError, "message": "Gagal membaca file"})
		return
	}

	// Connect to MySQL
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/hr_attendance_app")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer db.Close()

	query := "INSERT INTO gambar_user (image_path) VALUES (?)"
	_, err = db.Exec(query, imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menyimpan gambar ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gambar berhasil dikirim"})
}

func CompareImage(c *gin.Context, imageData []byte) {
	image, err := gocv.IMDecode(imageData, gocv.IMReadUnchanged)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Gagal mendecode gambarr"})
		return
	}

	defer image.Close()

	//proses matchmataking dengan gambar referensi di database
	matVal := gocv.NewMat()
	defer matVal.Close()

	gocv.CvtColor(image, &matVal, gocv.ColorBGR555ToGRAY)

	//proses matchmaking gambar
	newRef := gocv.IMRead("refrensi.jpg", gocv.IMReadGrayScale)
	if newRef.Empty() {
		fmt.Println("Error reading reference image")
		return
	}

	defer newRef.Close()

	res := gocv.NewMat()
	defer res.Close()

	gocv.MatchTemplate(image, newRef, &res, gocv.TmCcoeffNormed, res)

}
