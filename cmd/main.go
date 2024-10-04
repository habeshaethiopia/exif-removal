package main

import (
	"log"
	"os"

	"exif-removal/internal/controller"
	"exif-removal/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create necessary directories
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		err := os.Mkdir("uploads", os.ModePerm)
		if err != nil {
			log.Fatalf("failed to create uploads directory: %v", err)
		}
	}

	// Initialize services and controllers
	exifService := service.NewExifService()
	exifController := controller.NewExifController(exifService)

	// Set up Gin router
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	// Routes
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})
	router.POST("/upload", exifController.UploadHandler)

	// Start the server
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
