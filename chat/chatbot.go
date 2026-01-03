package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gin-gonic/gin"
)

func ChatBotOllama(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()
	var req modules.RequestChat

	if err := c.ShouldBindJSON(&req); err != nil || req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}

	//input post
	body, _ := json.Marshal(map[string]interface{}{
		"model":  "deepseek-r1:1.5b",
		"prompt": req.Prompt,
		"stream": false,
	})
	reqOllama, _ := http.NewRequestWithContext(ctx, os.Getenv("METHOD_CHAT"), os.Getenv("URL_CHAT"), bytes.NewReader(body))
	reqOllama.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqOllama)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghubungi model", "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result modules.RequestStreamChat
	if err := json.Unmarshal(bodyResp, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reply": result.Response,
	})

}
