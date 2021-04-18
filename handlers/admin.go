package handlers

import (
	"fmt"
	"go-drive/entities/admin"

	"github.com/gin-gonic/gin"
)

type adminHandler struct {
	adminService admin.Service
}

func NewAdminHandler(adminService admin.Service) *adminHandler {
	return &adminHandler{adminService: adminService}
}

func (h *adminHandler) ListRootDir(c *gin.Context) {
	files, err := h.adminService.ListRootDir()

	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": fmt.Sprintf("Error retrieving files: %v", err)})

		return
	}

	c.JSON(200, gin.H{"status": "success", "data": files})
}

func (h *adminHandler) ResetStorage(c *gin.Context) {
	files, err := h.adminService.ListRootDir()

	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": fmt.Sprintf("Error retrieving files: %v", err)})

		return
	}

	if len(files) >= 1 {
		if err = h.adminService.ResetStorage(files); err != nil {
			c.JSON(500, gin.H{"status": "fail", "message": fmt.Sprintf("Error retrieving files: %v", err)})

			return
		}
	}

	c.JSON(201, gin.H{"status": "deleted", "message": "Success resetting storage"})

}

func (h *adminHandler) GetAppDir(c *gin.Context) {
	driveFile, err := h.adminService.GetAppDir()

	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": fmt.Sprintf("Error retrieving files: %v", err)})

		return
	}

	if driveFile.ID == "" {
		c.JSON(404, gin.H{"status": "not-found", "message": "No app storage found"})

		return
	}

	c.JSON(200, gin.H{"status": "success", "data": driveFile})
}

func (h *adminHandler) CreateAppDir(c *gin.Context) {
	appDir, err := h.adminService.CreateAppDir()

	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": fmt.Sprintf("Error creating app directory: %v", err)})

		return
	}

	c.JSON(201, gin.H{"status": "created", "data": gin.H{"id": appDir.ID, "name": appDir.Name}})
}
