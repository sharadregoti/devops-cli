package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
)

var fClient *firestore.Client
var projectID = "try-out-gcp-features"

func IncrementAppStarts() error {
	if fClient == nil {
		return nil
	}

	ctx := context.Background()
	docRef := fClient.Collection("metrics").NewDoc()
	user := map[string]interface{}{
		"type":   "start_counter",
		"labels": getUniqueInfo(),
		"value":  1,
	}

	if _, err := docRef.Set(ctx, user); err != nil {
		log.Printf("Failed to set data: %v", err)
	}
	return nil
}

func getUniqueInfo() map[string]string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return map[string]string{}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return map[string]string{}
	}
	return map[string]string{
		"home": homeDir,
		"host": hostname,
	}
}

func ReportUsageTime(startTime, endTime time.Time) error {
	if fClient == nil {
		return nil
	}

	// Calculate the usage time
	usageTime := endTime.Sub(startTime)

	ctx := context.Background()
	docRef := fClient.Collection("metrics").NewDoc()
	user := map[string]interface{}{
		"type":   "usage_time",
		"labels": getUniqueInfo(),
		"value":  usageTime.Seconds(),
	}

	if _, err := docRef.Set(ctx, user); err != nil {
		log.Printf("Failed to set data: %v", err)
	}
	return nil
}
