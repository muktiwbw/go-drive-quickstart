package storage

import (
	"time"

	"google.golang.org/api/drive/v3"
)

type AccountService struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

type DriveFile struct {
	ID        string
	Name      string
	URL       string
	MimeType  string
	CreatedAt time.Time
}

func FormatDriveFile(f *drive.File) (DriveFile, error) {
	createdAt, err := time.Parse(time.RFC3339, f.CreatedTime)

	if err != nil {
		return DriveFile{}, err
	}

	return DriveFile{
		ID:        f.Id,
		Name:      f.Name,
		URL:       f.WebViewLink,
		MimeType:  f.MimeType,
		CreatedAt: createdAt,
	}, nil
}
