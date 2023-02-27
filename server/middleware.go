package server

// import (
// 	"bytes"
// 	"io"
// 	"net/http"

// 	"github.com/Azure/go-autorest/logger"
// )

// // LoggerMiddleWare logs all incomming request
// func LoggerMiddleWare(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var reqBody []byte
// 		if r.Header.Get("Content-Type") == "application/json" {
// 			reqBody, _ = io.ReadAll(r.Body)
// 			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
// 		}

// 		logger.LogInfo("Request", map[string]interface{}{"method": r.Method, "url": r.URL.Path, "queryVars": r.URL.Query(), "body": string(reqBody)})
// 		next.ServeHTTP(w, r)
// 	})
// }
