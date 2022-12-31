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

	// /// Prepares an individual data point
	// dataPoint := &monitoringpb.Point{
	// 	Interval: &monitoringpb.TimeInterval{
	// 		EndTime: &googlepb.Timestamp{
	// 			Seconds: time.Now().Unix(),
	// 		},
	// 	},
	// 	Value: &monitoringpb.TypedValue{
	// 		Value: &monitoringpb.TypedValue_Int64Value{
	// 			Int64Value: 1,
	// 		},
	// 	},
	// }

	// // Writes time series data.
	// // mClient.CallOptions
	// if err := mClient.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
	// 	Name: fmt.Sprintf("projects/%s", projectID),
	// 	TimeSeries: []*monitoringpb.TimeSeries{
	// 		{
	// 			Metric: &metricpb.Metric{
	// 				Type:   "custom.googleapis.com/app/start_counter",
	// 				Labels: getUniqueInfo(),
	// 			},
	// 			Resource: &monitoredrespb.MonitoredResource{
	// 				Type: "global",
	// 				Labels: map[string]string{
	// 					"project_id": projectID,
	// 				},
	// 			},
	// 			Points: []*monitoringpb.Point{
	// 				dataPoint,
	// 			},
	// 		},
	// 	},
	// }); err != nil {
	// 	log.Printf("Failed to write time series data: %v", err)
	// }
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
	// /// Prepares an individual data point
	// dataPoint := &monitoringpb.Point{
	// 	Interval: &monitoringpb.TimeInterval{
	// 		EndTime: &googlepb.Timestamp{
	// 			Seconds: time.Now().Unix(),
	// 		},
	// 	},
	// 	Value: &monitoringpb.TypedValue{
	// 		Value: &monitoringpb.TypedValue_Int64Value{
	// 			Int64Value: int64(usageTime.Seconds()),
	// 		},
	// 	},
	// }

	// // Writes time series data.
	// if err := mClient.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
	// 	Name: fmt.Sprintf("projects/%s", projectID),
	// 	TimeSeries: []*monitoringpb.TimeSeries{
	// 		{
	// 			Metric: &metricpb.Metric{
	// 				Type:   "custom.googleapis.com/app/usage_time",
	// 				Labels: getUniqueInfo(),
	// 			},
	// 			Resource: &monitoredrespb.MonitoredResource{
	// 				Type: "global",
	// 				Labels: map[string]string{
	// 					"project_id": projectID,
	// 				},
	// 			},
	// 			Points: []*monitoringpb.Point{
	// 				dataPoint,
	// 			},
	// 		},
	// 	},
	// }); err != nil {
	// 	log.Printf("Failed to write time series data: %v", err)
	// }
	return nil
}
