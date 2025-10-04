package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func TestRequestLogger(t *testing.T) {
	// Инициализация логгера для тестов
	err := logger.InitSimple("info", false)
	require.NoError(t, err)

	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments", nil)
	req = req.WithContext(context.Background())
	w := httptest.NewRecorder()

	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Act
	RequestLogger(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.True(t, handlerCalled, "Следующий handler должен быть вызван")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestLogger_POST(t *testing.T) {
	// Инициализация логгера для тестов
	err := logger.InitSimple("info", false)
	require.NoError(t, err)

	// Arrange
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", nil)
	req = req.WithContext(context.Background())
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	// Act
	RequestLogger(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRequestLogger_WithError(t *testing.T) {
	// Инициализация логгера для тестов
	err := logger.InitSimple("info", false)
	require.NoError(t, err)

	// Arrange
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/payments/123", nil)
	req = req.WithContext(context.Background())
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	// Act
	RequestLogger(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Act
	rw.WriteHeader(http.StatusCreated)

	// Assert
	assert.Equal(t, http.StatusCreated, rw.statusCode)
	assert.Equal(t, http.StatusCreated, w.Code)
}
