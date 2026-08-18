package middleware

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// TraceIDKey is the context key for trace ID.
type contextKey string

const traceIDKey contextKey = "traceid"

// TraceIDHeader is the HTTP header for trace ID.
const TraceIDHeader = "X-Trace-ID"

// Logger is the structured JSON logger.
type Logger struct {
	logger *log.Logger
}

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	TraceID   string `json:"trace_id,omitempty"`
	Message   string `json:"message"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	Status    int    `json:"status,omitempty"`
	Duration  string `json:"duration,omitempty"`
	Error     string `json:"error,omitempty"`
}

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(os.Stdout)
}

// DefaultLogger returns the default logger instance.
func DefaultLogger() *Logger {
	return defaultLogger
}

// MockLogger creates a logger that writes to an io.Writer (for testing).
func MockLogger(w io.Writer) *Logger {
	return NewLogger(w)
}

// NewLogger creates a new structured logger.
func NewLogger(output io.Writer) *Logger {
	return &Logger{
		logger: log.New(output, "", 0),
	}
}

func (l *Logger) log(entry logEntry) {
	data, _ := json.Marshal(entry)
	l.logger.Println(string(data))
}

// Info logs an info message.
func (l *Logger) Info(msg string, fields ...LogField) {
	entry := logEntry{Timestamp: time.Now().UTC().Format(time.RFC3339), Level: "INFO", Message: msg}
	for _, f := range fields {
		f(&entry)
	}
	l.log(entry)
}

// Error logs an error message.
func (l *Logger) Error(msg string, err error, fields ...LogField) {
	entry := logEntry{Timestamp: time.Now().UTC().Format(time.RFC3339), Level: "ERROR", Message: msg}
	if err != nil {
		entry.Error = err.Error()
	}
	for _, f := range fields {
		f(&entry)
	}
	l.log(entry)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, fields ...LogField) {
	entry := logEntry{Timestamp: time.Now().UTC().Format(time.RFC3339), Level: "WARN", Message: msg}
	for _, f := range fields {
		f(&entry)
	}
	l.log(entry)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, fields ...LogField) {
	entry := logEntry{Timestamp: time.Now().UTC().Format(time.RFC3339), Level: "DEBUG", Message: msg}
	for _, f := range fields {
		f(&entry)
	}
	l.log(entry)
}

// LogField is a field for structured logging.
type LogField func(*logEntry)

func TraceID(id string) LogField {
	return func(e *logEntry) {
		e.TraceID = id
	}
}

func Method(m string) LogField {
	return func(e *logEntry) {
		e.Method = m
	}
}

func Path(p string) LogField {
	return func(e *logEntry) {
		e.Path = p
	}
}

func Status(s int) LogField {
	return func(e *logEntry) {
		e.Status = s
	}
}

func Duration(d time.Duration) LogField {
	return func(e *logEntry) {
		e.Duration = d.String()
	}
}

// ExtractTraceID extracts trace ID from context.
func ExtractTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return ""
}

// NewTraceID generates a new trace ID.
func NewTraceID() string {
	return uuid.New().String()
}

// TraceIDMiddleware adds trace ID to context and response headers (chi-compatible).
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = NewTraceID()
		}

		ctx := context.WithValue(r.Context(), traceIDKey, traceID)
		w.Header().Set(TraceIDHeader, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestLogger is a middleware that logs HTTP requests in JSON format.
func RequestLogger(logger *Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Get trace ID from context if available
			traceID := ExtractTraceID(r.Context())

			// Wrap response writer to capture status
			wrapped := &statusWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			entry := logEntry{
				Timestamp: start.UTC().Format(time.RFC3339),
				Level:     "INFO",
				TraceID:   traceID,
				Method:    r.Method,
				Path:      r.URL.Path,
				Status:    wrapped.status,
				Duration:  duration.String(),
			}

			if wrapped.status >= 500 {
				entry.Level = "ERROR"
			} else if wrapped.status >= 400 {
				entry.Level = "WARN"
			}

			logger.log(entry)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}
