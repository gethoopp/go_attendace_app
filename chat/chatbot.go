package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gethoopp/hr_attendance_app/modules"
	"github.com/gin-gonic/gin"
)

func ChatBotOllama(c *gin.Context) {
	ctx := context.Background()
	var req modules.RequestChat

	if err := c.ShouldBindJSON(&req); err != nil || req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": gin.H{"message": "Gagal mengirimkan "}})
		return
	}

	//input post
	body, _ := json.Marshal(map[string]interface{}{
		"model":  "llama3.2",
		"prompt": req.Prompt,
	})
	reqOllama, _ := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/generate", bytes.NewReader(body))
	reqOllama.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqOllama)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghubungi model"})
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var reply string
	for scanner.Scan() {
		line := scanner.Bytes()
		var part modules.RequestStreamChat
		if err := json.Unmarshal(line, &part); err == nil {
			reply += part.Response
			if part.Done {
				break
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"reply": reply,
	})
}
