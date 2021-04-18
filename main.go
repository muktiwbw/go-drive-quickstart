package main

import (
	"encoding/json"
	"go-drive/entities/admin"
	"go-drive/entities/web"
	"go-drive/handlers"
	"go-drive/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// * Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env: %v\n", err)
	}

	// * Set project id env
	var accountService storage.AccountService

	if err := json.Unmarshal([]byte(os.Getenv("GOOGLE_ACCOUNT_SERVICE_JSON")), &accountService); err != nil {
		log.Fatalf("Error parsing account service JSON content: %v\n", err)
	}

	if err := os.Setenv("GOOGLE_PROJECT_ID", accountService.ProjectID); err != nil {
		log.Fatalf("Error setting project id: %v\n", err)
	}

	// * Init storage service using Google Drive API
	srv := storage.NewStorageService()

	// * Init services
	webService := web.NewWebService(srv)
	adminService := admin.NewService(srv)

	// * Init handlers
	webHandler := handlers.NewWebHandler(webService)
	adminHandler := handlers.NewAdminHandler(adminService)

	// * Init router engine
	router := gin.Default()
	api := router.Group("/api/v1")
	admin := router.Group("/api/v1/admin")

	// * Init routers
	api.GET("/", webHandler.GetStarterData)
	api.POST("/upload", webHandler.UploadAvatar)

	admin.GET("/storage", adminHandler.ListRootDir)
	admin.POST("/storage/reset", adminHandler.ResetStorage)
	admin.GET("/storage/appdir", adminHandler.GetAppDir)
	admin.POST("/storage/appdir", adminHandler.CreateAppDir)

	// TODO
	// ! Clean up files within app directory
	// ! Upload files
	// ! Delete files

	// * Run the server
	router.Run()
}
