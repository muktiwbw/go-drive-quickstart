package handlers

import (
	"fmt"
	"go-drive/entities/web"

	"github.com/gin-gonic/gin"
)

type webHandler struct {
	webService web.Service
}

func NewWebHandler(webService web.Service) *webHandler {
	return &webHandler{webService: webService}
}

func (h *webHandler) GetStarterData(c *gin.Context) {
	payload := gin.H{
		"message": "success",
		"data":    "Hi, cutie!",
	}

	c.JSON(200, payload)
}

func (h *webHandler) UploadAvatar(c *gin.Context) {
	avatar, err := c.FormFile("avatar")

	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": fmt.Sprintf("Error retrieving file: %v", err)})

		return
	}

	uploadedAvatar, err := h.webService.Upload(avatar)

	var statusCode int
	var payload map[string]interface{}

	if err != nil {
		statusCode = 500
		payload = gin.H{
			"status": "success",
			"data": gin.H{
				"file_id": uploadedAvatar.Id,
			},
		}
	} else {
		statusCode = 201
		payload = gin.H{
			"status": "success",
			"data": gin.H{
				"file_id": uploadedAvatar.Id,
			},
		}
	}

	c.JSON(statusCode, payload)
}
