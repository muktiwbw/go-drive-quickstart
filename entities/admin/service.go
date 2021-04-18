package admin

import (
	"fmt"
	"go-drive/storage"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
)

type Service interface {
	ListRootDir() ([]storage.DriveFile, error)
	ResetStorage(files []storage.DriveFile) error
	GetAppDir() (storage.DriveFile, error)
	CreateAppDir() (storage.DriveFile, error)
	CleanStorageApp(dirID string) error
}

type service struct {
	driveService *drive.Service
}

func NewService(driveService *drive.Service) Service {
	return &service{driveService: driveService}
}

func (s *service) ListRootDir() ([]storage.DriveFile, error) {
	fileList := []storage.DriveFile{}

	driveFileList, err := s.driveService.Files.List().Fields("files(id, name, webViewLink, mimeType, createdTime)").Do()

	if err != nil {
		return fileList, err
	}

	for _, file := range driveFileList.Files {
		formattedDriveFile, err := storage.FormatDriveFile(file)

		if err != nil {
			return fileList, err
		}

		fileList = append(fileList, formattedDriveFile)
	}

	return fileList, nil

}

func (s *service) ResetStorage(files []storage.DriveFile) (err error) {
	for _, file := range files {
		err = s.driveService.Files.Delete(file.ID).Do()

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) GetAppDir() (storage.DriveFile, error) {
	appName := fmt.Sprintf("storage_%s_%s", os.Getenv("GOOGLE_PROJECT_ID"), os.Getenv("APP_NAME"))
	appDirID := os.Getenv("DRIVE_APP_DIR_ID")
	var appDir *drive.File

	// * If you provides appdir id, check by the id
	if appDirID != "" {
		foundAppDir, err := s.driveService.Files.Get(appDirID).Fields("id, name, webViewLink, mimeType, createdTime").Do()

		// * Check if err is caused by something other than not found
		if err != nil {
			e := strings.Split(err.Error(), ", ")
			if e[len(e)-1] != "notFound" {
				return storage.DriveFile{}, err
			}
		}

		// * Assign to appDir if file exists
		// * at this point the file must definitely be available
		if err == nil {
			appDir = foundAppDir
		}
	}

	// * If appDir is still nil, which means file not found after searching by id, search it by the app name
	if appDir == nil {
		driveQuery := fmt.Sprintf("(name='%s' and mimeType='application/vnd.google-apps.folder')", appName)

		fileList, err := s.driveService.Files.List().Q(driveQuery).Fields("files(id, name, webViewLink, mimeType, createdTime)").Do()

		if err != nil {
			return storage.DriveFile{}, err
		}

		if len(fileList.Files) <= 0 {
			return storage.DriveFile{}, nil
		}

		appDir = fileList.Files[0]
	}

	formattedDriveFile, err := storage.FormatDriveFile(appDir)

	if err != nil {
		return formattedDriveFile, err
	}

	return formattedDriveFile, nil

}

func (s *service) CreateAppDir() (storage.DriveFile, error) {
	appName := fmt.Sprintf("storage_%s_%s", os.Getenv("GOOGLE_PROJECT_ID"), os.Getenv("APP_NAME"))

	appDir, err := s.driveService.Files.Create(&drive.File{Name: appName, MimeType: "application/vnd.google-apps.folder"}).Do()

	if err != nil {
		return storage.DriveFile{}, err
	}

	// * Set to read only permission for anyone
	_, err = s.driveService.Permissions.Create(appDir.Id, &drive.Permission{Type: "anyone", Role: "reader"}).Do()

	if err != nil {
		return storage.DriveFile{}, err
	}

	// ? Adding your real email so that you can easily organize sub-folders in the website
	organizerEmail := os.Getenv("DRIVE_ORGANIZER_EMAIL")
	if organizerEmail != "" {
		_, err = s.driveService.Permissions.Create(appDir.Id, &drive.Permission{Type: "user", Role: "writer", EmailAddress: organizerEmail}).Do()

		if err != nil {
			return storage.DriveFile{}, err
		}
	}

	return storage.DriveFile{ID: appDir.Id, Name: appDir.Name}, nil
}

func (s *service) CleanStorageApp(dirID string) error {

	return nil
}
