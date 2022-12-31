package common

import (
	"context"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"google.golang.org/api/option"

	"log"
)

var data = `
{
	"client_id": "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com",
	"client_secret": "d-FL95Q19q7MQmFpd7hHD0Ty",
	"quota_project_id": "try-out-gcp-features",
	"refresh_token": "1//0gEi9-8AIjfiVCgYIARAAGBASNwF-L9Ir0stzEqkcB-y0MLsvg9DoBW_8o2fzXeYF9a5Zir-1VL9QXz-vjZiH89OsQ2kcPrdBdSs",
	"type": "authorized_user"
}`

func ConnInit() {

	ctx := context.Background()

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(data)))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}
	// defer client.Close()

	// Sets the name of the log to write to.
	logName := "devops-cli"

	logger = client.Logger(logName).StandardLogger(logging.Error)

	clientm, err := firestore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(data)))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	fClient = clientm
}

func ConnLoggingInit() {
	ctx := context.Background()

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(data)))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}
	// defer client.Close()

	// Sets the name of the log to write to.
	logName := "devops-cli"

	logger = client.Logger(logName).StandardLogger(logging.Error)
}
