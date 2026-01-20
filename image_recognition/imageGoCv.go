package imageRecognition

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gethoopp/hr_attendance_app/database"
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
	db, err := database.GetDB()
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

func CompareImage(c *gin.Context, imageData []byte, imageDB []byte) {

	image, err := gocv.IMDecode(imageData, gocv.IMReadUnchanged)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Gagal mendecode gambar"})
		return
	}
	defer image.Close()

	matVal := gocv.NewMat()
	defer matVal.Close()
	gocv.CvtColor(image, &matVal, gocv.ColorBGRToGray)

	refImage, err := gocv.IMDecode(imageDB, gocv.IMReadGrayScale)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Error decoding reference image"})
		return
	}
	if refImage.Empty() {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Error reading reference image"})
		return
	}
	defer refImage.Close()

	res := gocv.NewMat()
	defer res.Close()

	// Perform template matching
	gocv.MatchTemplate(matVal, refImage, &res, gocv.TmCcoeffNormed, refImage)

	minVal, maxVal, minLoc, maxLoc := gocv.MinMaxLoc(res)
	threshold := float32(0.8)
	if maxVal >= threshold {
		c.JSON(http.StatusOK, gin.H{
			"status":     http.StatusOK,
			"message":    "Image match found",
			"location":   maxLoc,
			"confidence": maxVal,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":     http.StatusOK,
			"message":    "No significant match found",
			"location":   minLoc,
			"confidence": minVal,
		})
	}
}
