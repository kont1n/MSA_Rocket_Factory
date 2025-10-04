package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Handle(t *testing.T) {
	// Arrange
	handler := NewHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Act
	handler.Handle(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"status":"ok","service":"assembly"}`, w.Body.String())
}

func TestHealthHandler_Handle_MultipleRequests(t *testing.T) {
	// Arrange
	handler := NewHealthHandler()

	tests := []struct {
		name   string
		method string
	}{
		{
			name:   "GET запрос",
			method: http.MethodGet,
		},
		{
			name:   "POST запрос",
			method: http.MethodPost,
		},
		{
			name:   "HEAD запрос",
			method: http.MethodHead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			// Act
			handler.Handle(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tt.method != http.MethodHead {
				assert.JSONEq(t, `{"status":"ok","service":"assembly"}`, w.Body.String())
			}
		})
	}
}

func TestHealthHandler_Handle_ConcurrentRequests(t *testing.T) {
	// Arrange
	handler := NewHealthHandler()
	const numRequests = 10

	// Act - выполняем несколько запросов конкурентно
	done := make(chan bool, numRequests)
	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			handler.Handle(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			done <- true
		}()
	}

	// Assert - ждем завершения всех запросов
	for i := 0; i < numRequests; i++ {
		<-done
	}
}

func TestNewHealthHandler(t *testing.T) {
	// Act
	handler := NewHealthHandler()

	// Assert
	assert.NotNil(t, handler)
}

func TestHealthHandler_Handle_ResponseBody(t *testing.T) {
	// Arrange
	handler := NewHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Act
	handler.Handle(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "status")
	assert.Contains(t, body, "ok")
	assert.Contains(t, body, "service")
	assert.Contains(t, body, "assembly")
}
