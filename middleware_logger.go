package az

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(MethodStruct string) string {
	switch MethodStruct {
	case "GET":
		return blue
	case "POST":
		return green
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return cyan
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

func logger(h http.HandlerFunc, _ ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		h(lrw, r)
		end := time.Now()
		latency := end.Sub(start)

		clientIP := r.RemoteAddr
		MethodStruct := r.Method
		statusCode := lrw.statusCode
		var statusColor, methodColor string
		statusColor = colorForStatus(statusCode)
		methodColor = colorForMethod(MethodStruct)

		fmt.Printf("[AZ] %v |%s %3d %s| %13v | %s |%s %-7s %s| %s | %s\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, MethodStruct, reset,
			r.URL.Path,
			strings.Replace(r.URL.RawQuery, "&", " ", -1),
		)
	}
}
