package service

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"

	"github.com/rwcarlsen/goexif/exif"
)

type ExifService interface {
	CheckExif(file io.Reader) (bool, error)
	RemoveExif(file io.Reader) ([]byte, error)
}

type exifService struct{}

func NewExifService() ExifService {
	return &exifService{}
}

// CheckExif verifies if EXIF data is present in the image
func (s *exifService) CheckExif(file io.Reader) (bool, error) {
	_, err := exif.Decode(file)
	if err == nil {
		return true, nil // EXIF data found
	}
	return false, fmt.Errorf("no EXIF metadata found: %w", err)
}

// RemoveExif strips EXIF metadata from the image
func (s *exifService) RemoveExif(file io.Reader) ([]byte, error) {
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, img, nil); err != nil {
		return nil, fmt.Errorf("failed to encode image without EXIF: %w", err)
	}

	return buffer.Bytes(), nil
}
