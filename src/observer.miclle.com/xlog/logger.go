package xlog

import (
	"log"
	"net/http"
	"os"

	"github.com/satori/go.uuid"
)

// A Logger represents an active logging object
type Logger struct {
	*log.Logger
	header     http.Header
	xRequestID string
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{
		Logger:     log.New(os.Stdout, "[observer] ", 0),
		xRequestID: uuid.NewV4().String(),
	}
}

// XRequestID return X-Request-ID
func (xlog *Logger) XRequestID() string {
	return xlog.xRequestID
}
