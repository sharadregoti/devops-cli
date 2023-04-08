package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	// "github.com/go-errors/errors"
	"github.com/rs/cors"
	"github.com/sharadregoti/devops/model"
	// "github.com/xebiaww-apps/xlr8s-go/model/httptypes"
	// "github.com/xebiaww-apps/xlr8s-go/utils/errors"
)

// CreateCorsObject creates a cors object with the required config
func CreateCorsObject() *cors.Cors {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc: func(s string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		ExposedHeaders: []string{"Authorization", "Content-Type"},
	})
}

// CloseTheCloser closes the closer
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}

// SendResponse sends an http response
func SendResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) error {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(body)
}

// SendErrorResponse sends an http error response
func SendErrorResponse(ctx context.Context, w http.ResponseWriter, statusCode int, err error) error {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	fmt.Println("Errror resposne 1 status has been sent")
	response := model.ErrorResponse{
		Message: err.Error(),
	}

	fmt.Println("Errror resposne 2 message bing sent")
	return json.NewEncoder(w).Encode(response)
}

// ExtractToken extracts token from http request
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	} else if len(strArr) == 1 {
		return strArr[0]
	}
	return ""
}
