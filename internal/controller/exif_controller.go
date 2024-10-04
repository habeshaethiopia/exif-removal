package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"exif-removal/internal/service"

	"github.com/gin-gonic/gin"
)

type ExifController struct {
	ExifService service.ExifService
}

func NewExifController(svc service.ExifService) *ExifController {
	return &ExifController{ExifService: svc}
}

func (c *ExifController) UploadHandler(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read file"})
		return
	}
	defer file.Close()

	// Check for EXIF metadata
	containsExif, err := c.ExifService.CheckExif(file)
	if err != nil && !containsExif {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "File does not contain EXIF metadata or an error occurred"})
		return
	}

	// Reset the file pointer to the start after reading for EXIF
	file.Seek(0, 0)

	// Remove EXIF metadata
	cleanedData, err := c.ExifService.RemoveExif(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove EXIF metadata"})
		return
	}

	// Save the cleaned image
	cleanedFilePath := filepath.Join("uploads", "cleaned_"+header.Filename)
	err = os.WriteFile(cleanedFilePath, cleanedData, 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save cleaned image"})
		return
	}

	// Respond with the success message
	ctx.JSON(http.StatusOK, gin.H{
		"message":           "EXIF metadata removed successfully",
		"cleaned_image_url": cleanedFilePath,
	})
}
