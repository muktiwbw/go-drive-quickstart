package storage

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func NewStorageService() *drive.Service {
	serviceAccountJsonContent := os.Getenv("GOOGLE_ACCOUNT_SERVICE_JSON")

	// * Generate a JSON file based on your service account data stored in env
	if _, err := os.Stat("./svracc.json"); err != nil && os.IsNotExist(err) {
		err := ioutil.WriteFile("./svracc.json", []byte(serviceAccountJsonContent), os.ModePerm)

		if err != nil {
			log.Fatalf("Error writing service account file: %v", err)
		}
	} else if err == nil {
		// * If file actually exists
		jsonContent, err := os.ReadFile("./svracc.json")
		if err != nil {
			log.Fatalf("Error reading service account file: %v", err)
		}

		// * Check if the content is the same as env or not
		// * If false, override with the env content
		if string(jsonContent) != serviceAccountJsonContent {
			err := ioutil.WriteFile("./svracc.json", []byte(serviceAccountJsonContent), os.ModePerm)

			if err != nil {
				log.Fatalf("Error overriding service account file: %v", err)
			}
		}
	} else if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading service account file: %v", err)
	}

	// * Create a new drive service
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsFile("./svracc.json"))

	if err != nil {
		log.Fatalf("Error creating drive service api: %v", err)
	}

	return srv
}
