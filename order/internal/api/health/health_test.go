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
	assert.JSONEq(t, `{"status":"ok","service":"order"}`, w.Body.String())
}

func TestNewHealthHandler(t *testing.T) {
	// Act
	handler := NewHealthHandler()

	// Assert
	assert.NotNil(t, handler)
}
