package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPMetrics(t *testing.T) {
	// Act
	metrics, err := NewHTTPMetrics()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.requestCount)
	assert.NotNil(t, metrics.requestDuration)
	assert.NotNil(t, metrics.errorCount)
}

func TestMetricsMiddleware_Success(t *testing.T) {
	// Arrange
	metrics, err := NewHTTPMetrics()
	require.NoError(t, err)

	middleware := MetricsMiddleware(metrics)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments", nil)
	w := httptest.NewRecorder()

	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Act
	middleware(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.True(t, handlerCalled, "Следующий handler должен быть вызван")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMetricsMiddleware_ErrorStatus(t *testing.T) {
	// Arrange
	metrics, err := NewHTTPMetrics()
	require.NoError(t, err)

	middleware := MetricsMiddleware(metrics)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	// Act
	middleware(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMetricsMiddleware_WithUUID(t *testing.T) {
	// Arrange
	metrics, err := NewHTTPMetrics()
	require.NoError(t, err)

	middleware := MetricsMiddleware(metrics)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Act
	middleware(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMetricsResponseWriter_WriteHeader(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	rw := &metricsResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Act
	rw.WriteHeader(http.StatusAccepted)

	// Assert
	assert.Equal(t, http.StatusAccepted, rw.statusCode)
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "health endpoint",
			path:     "/health",
			expected: "/health",
		},
		{
			name:     "payments list",
			path:     "/api/v1/payments",
			expected: "/api/v1/payments",
		},
		{
			name:     "payment by UUID",
			path:     "/api/v1/payments/550e8400-e29b-41d4-a716-446655440000",
			expected: "/api/v1/payments/{uuid}",
		},
		{
			name:     "unknown path",
			path:     "/api/v2/test",
			expected: "/api/v2/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := normalizePath(tt.path)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}
