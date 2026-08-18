package handlers

import (
	"encoding/json"
	"net/http"

	"competitive-analysis-agent/internal/api/middleware"
)

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}, traceID string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(middleware.TraceIDHeader, traceID)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log but don't fail if response already started
		middleware.DefaultLogger().Error("Failed to encode JSON response", err)
	}
}

// writeError writes an error JSON response.
func writeError(w http.ResponseWriter, status int, message, traceID string) {
	type errorResponse struct {
		Error   string `json:"error"`
		TraceID string `json:"trace_id"`
	}
	writeJSON(w, status, errorResponse{Error: message, TraceID: traceID}, traceID)
}
