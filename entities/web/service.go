package web

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/api/drive/v3"
)

type Service interface {
	Upload(file *multipart.FileHeader) (*drive.File, error)
}

type service struct {
	driveService *drive.Service
}

func NewWebService(srv *drive.Service) Service {
	return &service{driveService: srv}
}

func (s *service) Upload(file *multipart.FileHeader) (*drive.File, error) {
	var driveFile *drive.File

	src, err := file.Open()
	if err != nil {
		return driveFile, err
	}
	defer src.Close()

	fileName := fmt.Sprintf("%s_%d", strings.Split(file.Filename, ".")[0], time.Now().UnixNano())
	fileExt := filepath.Ext(file.Filename)
	fullFileName := fileName + fileExt
	mimeType := "image/jpeg"

	driveFile, err = s.driveService.Files.Create(&drive.File{Name: fullFileName, MimeType: mimeType}).Media(src).Do()

	if err != nil {
		return driveFile, err
	}

	// ! Later change role to reader
	_, err = s.driveService.Permissions.Create(driveFile.Id, &drive.Permission{Type: "anyone", Role: "reader"}).Do()

	if err != nil {
		return driveFile, err
	}

	return driveFile, nil
}
