package httpheader

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

// SetResponseHeader sets HTTP response header
func SetResponseHeader(ctx *gin.Context) {
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.Writer.Header().Set("X-XSS-Protection", "1, mode=block")
	ctx.Writer.Header().Set("X-Frame-Options", "deny")
	ctx.Writer.Header().Set("Content-Security-Policy", "default-src 'none'")
	// ctx.Writer.Header().Set("Strict-Transport-Security", "max-age=15768000")
}

// setResponseHeader sets HTTP response header
//func setResponseHeader(ctx *gin.Context, data []map[string]string) {
//	for _, header := range data {
//		for key, val := range header {
//			ctx.Writer.Header().Set(key, val)
//		}
//	}
//}

// SetHTTPHeaders sets http request header
func SetHTTPHeaders(req *http.Request, headers []map[string]string) {
	// req.Header.Set("Authorization", "Bearer access-token")
	for _, header := range headers {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
}

// DebugHTTPHeader debugs http request headers
func DebugHTTPHeader(req *http.Request) error {
	dumped, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	log.Printf("[DebugHTTPHeader] headers:\n%s\n", dumped)
	// POST /login HTTP/1.1
	// Host: 127.0.0.1:63513
	// User-Agent: Go-http-client/1.1
	// Content-Length: 0
	// Content-Type: application/x-www-form-urlencoded
	// Cookie: go-web-ginserver=MTQ3MTA1MDQ3MnxOd3dBTkVOQlJGZE1WRTlRVmxoWldGbEVSVTFYVGxKSk5VZFhXalJYVkRWRlNWazJWRnBQVUVWWlJGSklOMUZSUTB0TE0waGFRVkU9fC_7LJ1pOXIOZo8ZXg-R4oO1LFXaSqJtvA3l0f6Qk9DA
	// Referer: http://hiromaily.com:8080/login
	// Accept-Encoding: gzip

	return nil
}
