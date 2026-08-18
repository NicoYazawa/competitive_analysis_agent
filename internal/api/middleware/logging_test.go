package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Info("test message",
		TraceID("trace-123"),
		Method("GET"),
		Path("/api/test"),
	)

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "INFO", entry.Level)
	assert.Equal(t, "test message", entry.Message)
	assert.Equal(t, "trace-123", entry.TraceID)
	assert.Equal(t, "GET", entry.Method)
	assert.Equal(t, "/api/test", entry.Path)
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Error("error occurred", assert.AnError,
		TraceID("trace-456"),
	)

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "ERROR", entry.Level)
	assert.Equal(t, "error occurred", entry.Message)
	assert.Equal(t, "trace-456", entry.TraceID)
	assert.Contains(t, entry.Error, "assert.AnError")
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Warn("warning message", Path("/api/warn"))

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "WARN", entry.Level)
	assert.Equal(t, "warning message", entry.Message)
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Debug("debug message")

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "DEBUG", entry.Level)
	assert.Equal(t, "debug message", entry.Message)
}

func TestDefaultLogger_ReturnsNonNil(t *testing.T) {
	logger := DefaultLogger()
	assert.NotNil(t, logger)
}

func TestMockLogger_ReturnsNonNil(t *testing.T) {
	var buf bytes.Buffer
	logger := MockLogger(&buf)
	assert.NotNil(t, logger)
}

func TestMockLogger_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := MockLogger(&buf)

	logger.Info("test message")

	assert.NotEmpty(t, buf.String())
	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)
	assert.Equal(t, "INFO", entry.Level)
	assert.Equal(t, "test message", entry.Message)
}

func TestLogger_Debug_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Debug("debug with fields",
		TraceID("debug-trace"),
		Method("DEBUG"),
		Path("/api/debug"),
		Status(200),
		Duration(50*time.Millisecond),
	)

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "DEBUG", entry.Level)
	assert.Equal(t, "debug with fields", entry.Message)
	assert.Equal(t, "debug-trace", entry.TraceID)
	assert.Equal(t, "DEBUG", entry.Method)
	assert.Equal(t, "/api/debug", entry.Path)
	assert.Equal(t, 200, entry.Status)
	assert.Equal(t, "50ms", entry.Duration)
}

func TestLogger_Debug_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Debug("")

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "DEBUG", entry.Level)
	assert.Empty(t, entry.Message)
}

func TestLogFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Info("test",
		TraceID("tid"),
		Method("POST"),
		Path("/path"),
		Status(200),
		Duration(100*time.Millisecond),
	)

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "tid", entry.TraceID)
	assert.Equal(t, "POST", entry.Method)
	assert.Equal(t, "/path", entry.Path)
	assert.Equal(t, 200, entry.Status)
	assert.Equal(t, "100ms", entry.Duration)
}

func TestExtractTraceID(t *testing.T) {
	ctx := context.WithValue(context.Background(), traceIDKey, "test-trace-id")
	assert.Equal(t, "test-trace-id", ExtractTraceID(ctx))
	assert.Equal(t, "", ExtractTraceID(context.Background()))
}

func TestNewTraceID(t *testing.T) {
	id1 := NewTraceID()
	id2 := NewTraceID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.Len(t, id1, 36) // UUID format
}

func TestTraceIDMiddleware(t *testing.T) {
	var traceID string
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID = ExtractTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	handler := TraceIDMiddleware(nextHandler)

	// Test without trace ID header - should generate new one
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, traceID)
	assert.Equal(t, w.Header().Get(TraceIDHeader), traceID)

	// Test with trace ID header - should use existing one
	req2 := httptest.NewRequest("GET", "/test2", nil)
	req2.Header.Set(TraceIDHeader, "custom-trace-id")
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	assert.Equal(t, "custom-trace-id", traceID)
	assert.Equal(t, "custom-trace-id", w2.Header().Get(TraceIDHeader))
}

func TestRequestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestLogger(logger)(nextHandler)

	req := httptest.NewRequest("GET", "/api/test?param=value", nil)
	req.Header.Set(TraceIDHeader, "req-trace-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, buf.String())

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "INFO", entry.Level)
	assert.Equal(t, "GET", entry.Method)
	assert.Equal(t, "/api/test", entry.Path)
	assert.Equal(t, 200, entry.Status)
	// Note: trace_id is empty because TraceIDMiddleware sets it in context,
	// but RequestLogger only reads from context (which doesn't have it without middleware)
	assert.NotEmpty(t, entry.Duration)
}

func TestRequestLogger_ErrorStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	handler := RequestLogger(logger)(nextHandler)

	req := httptest.NewRequest("GET", "/api/error", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "ERROR", entry.Level)
	assert.Equal(t, 500, entry.Status)
}

func TestRequestLogger_WarnStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	handler := RequestLogger(logger)(nextHandler)

	req := httptest.NewRequest("GET", "/api/notfound", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "WARN", entry.Level)
	assert.Equal(t, 404, entry.Status)
}

func TestLogFieldNoOp(t *testing.T) {
	// Test that empty LogFields don't panic
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	logger.Info("test message") // No fields

	var entry logEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "INFO", entry.Level)
	assert.Equal(t, "test message", entry.Message)
	assert.Empty(t, entry.TraceID)
	assert.Empty(t, entry.Method)
	assert.Empty(t, entry.Path)
	assert.Equal(t, 0, entry.Status)
	assert.Empty(t, entry.Duration)
}

func TestUUIDGeneration(t *testing.T) {
	// Verify UUID format
	id := uuid.New()
	assert.NotEmpty(t, id.String())
	assert.Len(t, id.String(), 36)
}
