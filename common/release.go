//go:build release
// +build release

package common

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"context"
	"google.golang.org/api/option"

	"log"
)

func init() {
	Release = true

	data := `
{
	"client_id": "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com",
	"client_secret": "d-FL95Q19q7MQmFpd7hHD0Ty",
	"quota_project_id": "try-out-gcp-features",
	"refresh_token": "1//0gEi9-8AIjfiVCgYIARAAGBASNwF-L9Ir0stzEqkcB-y0MLsvg9DoBW_8o2fzXeYF9a5Zir-1VL9QXz-vjZiH89OsQ2kcPrdBdSs",
	"type": "authorized_user"
}`

	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	// projectID := "try-out-gcp-features"

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(data)))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}
	// defer client.Close()

	// Sets the name of the log to write to.
	logName := "devops-cli"

	logger = client.Logger(logName).StandardLogger(logging.Error)

	// Creates a client.
	// clientm, err := monitoring.NewMetricClient(ctx, option.WithCredentialsJSON([]byte(data)))
	// if err != nil {
	// 	log.Printf("Failed to create client: %v", err)
	// }
	clientm, err := firestore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(data)))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	fClient = clientm

	// Logs "hello world", log entry is visible at
	// Cloud Logs.
}
