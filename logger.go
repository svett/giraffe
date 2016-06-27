package giraffe

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// ColorGreen is a green color
	ColorGreen = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	// ColorWhite is a white color
	ColorWhite = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	// ColorYellow is a yellow color
	ColorYellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	// ColorRed is a red color
	ColorRed = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	// ColorBlue is a blue color
	ColorBlue = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	// ColorMagenta is a magenta color
	ColorMagenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	// ColorCyan is a cyan color
	ColorCyan = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	// DefaultColor is a default color
	DefaultColor = string([]byte{27, 91, 48, 109})
)

//go:generate counterfeiter -o fakes/fake_logger.go . Logger

// Logger that logs information
type Logger interface {
	// Println writes an info message
	Println(...interface{})
}

// HandlerFunc is a func that handle middleware operations
type HandlerFunc func(w http.ResponseWriter, request *http.Request, next http.HandlerFunc)

// NewHTTPStandardLogger prints logs into the standard out
func NewHTTPStandardLogger() HandlerFunc {
	color := false
	if info, err := os.Stderr.Stat(); err == nil {
		color = (info.Mode()&os.ModeCharDevice != 0)
	}
	return NewHTTPLogger(log.New(os.Stdout, "HTTP ", log.LstdFlags), color)
}

// NewHTTPLogger logs a HTTP requests
func NewHTTPLogger(logger Logger, color bool) HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		// Start timer
		start := time.Now()
		path := request.URL.Path

		// Process request
		writer := &responseWriter{ResponseWriter: w}
		next(writer, request)

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := request.RemoteAddr
		method := request.Method
		statusCode := writer.Status()

		var (
			statusColor string
			methodColor string
			resetColor  string
		)

		if color {
			statusColor = colorForStatus(statusCode)
			methodColor = colorForMethod(method)
			resetColor = DefaultColor
		}

		msg := fmt.Sprintf("%s %3d %s| %13v | %s |%s  %s %-7s %s",
			statusColor, statusCode, resetColor,
			latency,
			clientIP,
			methodColor, resetColor, method,
			path,
		)

		logger.Println(msg)
	}
}

// HTTPResponseWriter writes an response
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return ColorGreen
	case code >= 300 && code < 400:
		return ColorWhite
	case code >= 400 && code < 500:
		return ColorYellow
	default:
		return ColorRed
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return ColorBlue
	case "POST":
		return ColorCyan
	case "PUT":
		return ColorYellow
	case "DELETE":
		return ColorRed
	case "PATCH":
		return ColorGreen
	case "HEAD":
		return ColorMagenta
	case "OPTIONS":
		return ColorWhite
	default:
		return DefaultColor
	}
}
