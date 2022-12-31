package common

import (
	"log"

	"github.com/hashicorp/go-hclog"
)

var Release = false

var logger *log.Logger

// Send error to gcloud when using release binary
func Error(l hclog.Logger, msg string) {
	if Release && logger != nil {
		logger.Println(msg)
	}
	l.Error(msg)
}
