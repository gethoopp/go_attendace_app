package imageRecognition

import (
	"fmt"
	"image"
	"io"
	"net/http"

	"github.com/gethoopp/hr_attendance_app/database"
	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
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

	userID, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
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

	img, err := gocv.IMDecode(imageData, gocv.IMReadColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": http.StatusInternalServerError, "message": "Gagal konversi file"})
		return
	}

	net := gocv.ReadNetFromCaffe("deploy.prototxt", "res10_300x300_ssd_iter_140000_fp16.caffemodel")
	defer net.Close()

	blob := gocv.BlobFromImage(img, 1.0, image.Pt(300, 300), gocv.NewScalar(104, 177, 123, 0), false, false)
	defer blob.Close()
	net.SetInput(blob, "data")

	//detection
	detBlob := net.Forward("detection_out")
	defer detBlob.Close()

	detection := gocv.GetBlobChannel(detBlob, 0, 0)
	defer detection.Close()

	//ambil deteksi

	if detection.Rows() > 0 {
		confidence := detection.GetFloatAt(0, 2)

		if confidence > 0.5 {
			left := detection.GetFloatAt(0, 3) * float32(img.Cols())
			top := detection.GetFloatAt(0, 4) * float32(img.Rows())
			right := detection.GetFloatAt(0, 5) * float32(img.Cols())
			bottom := detection.GetFloatAt(0, 6) * float32(img.Rows())

			//cropping
			rect := image.Rect(int(left), int(top), int(right), int(bottom))
			faceImg := img.Region(rect)
			defer faceImg.Close()

			//hashing

			hashAlgo := contrib.PHash{}
			hasMat := gocv.NewMat()
			defer hasMat.Close()
			hashAlgo.Compute(faceImg, &hasMat)

			finalHash := hasMat.ToBytes()

			//get db
			db, err := database.GetDB()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
				return
			}
			defer db.Close()

			query := "UPDATE attendances SET image_hash = ? WHERE user_id = ?"
			_, err = db.Exec(query, finalHash, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Gagal mengupdate gambar ke database",
				})
				return
			}

		}
	}

	defer img.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Gambar berhasil dikirim"})
}

// func Compare_image_user(c *gin.Context) {
// 	var log modules.DataImage
// 	ctx, cancel := context.WithTimeout(c, 200*time.Second)
// 	defer cancel()

// 	imageResult, err := c.FormFile("image")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Gagal Mengirimkan gambar ke server"})
// 		return
// 	}

// 	file, err := imageResult.Open()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Status": http.StatusInternalServerError, "message": "Gagal membuka file"})
// 		return
// 	}
// 	defer file.Close()

// 	imageData, err := io.ReadAll(file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Status": http.StatusInternalServerError, "message": "Gagal membaca file"})
// 		return
// 	}

// 	userID, exists := c.Get("id_user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Unauthorized",
// 		})
// 		return
// 	}

// 	//koneksi db

// 	db, err := database.GetDB()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
// 		return
// 	}

// 	queries := "Select image_user from attendances where id = ?"
// 	rows, err := db.QueryContext(ctx, queries, userID)

// 	if err := rows.Scan(
// 		&log.ImageRef,
// 	); err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, gin.H{
// 				"status":  http.StatusNotFound,
// 				"message": "data user tidak ditemukan",
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  http.StatusBadRequest,
// 			"message": "Gagal membaca data",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	errs := CompareImage(c, imageData, log.ImageRef)
// 	if errs != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menyimpan gambar ke database"})
// 		return
// 	}

// 	defer db.Close()

// }

func CompareImage(c *gin.Context, imageData []byte, imageDB []byte) error {

	image, err := gocv.IMDecode(imageData, gocv.IMReadUnchanged)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Gagal mendecode gambar"})
		return nil
	}
	defer image.Close()

	matVal := gocv.NewMat()
	defer matVal.Close()
	gocv.CvtColor(image, &matVal, gocv.ColorBGRToGray)

	refImage, err := gocv.IMDecode(imageDB, gocv.IMReadGrayScale)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Error decoding reference image"})
		return err
	}
	if refImage.Empty() {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Error reading reference image"})
		return err
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

	return err
}
