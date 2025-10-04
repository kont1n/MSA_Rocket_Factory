package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
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

	// Проверяем, что логирование произошло
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Начало запроса")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/api/v1/orders")
	assert.Contains(t, logOutput, "Запрос завершен")
}

func TestRequestLogger_POST(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	// Act
	RequestLogger(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "POST")
	assert.Contains(t, logOutput, "/api/v1/orders")
}

func TestRequestLogger_WithError(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/123", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	// Act
	RequestLogger(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "DELETE")
	assert.Contains(t, logOutput, "Запрос завершен")
}
