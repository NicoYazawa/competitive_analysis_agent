package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.NotNil(t, router)
}

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
	assert.Contains(t, w.Body.String(), "ok")
}

func TestAPIv1ProductsNotFound(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIv1ProductsWithIDNotFound(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIv1CompetitorsNotFound(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIv1TrendsNotFound(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIv1StrategyPricingNotFound(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNotFoundHandler(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent/path", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Chi returns 405 for method not allowed on registered routes
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
