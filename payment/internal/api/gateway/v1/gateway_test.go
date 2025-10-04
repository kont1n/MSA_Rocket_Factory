package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
)

func TestNewGateway(t *testing.T) {
	// Act
	gateway := NewGateway()

	// Assert
	assert.NotNil(t, gateway)
	assert.NotNil(t, gateway.mux)
	assert.Nil(t, gateway.server)
}

func TestGateway_GetMux(t *testing.T) {
	// Arrange
	gateway := NewGateway()

	// Act
	mux := gateway.GetMux()

	// Assert
	assert.NotNil(t, mux)
	assert.Equal(t, gateway.mux, mux)
}

func TestGateway_SetHandler_WithServeMux(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	newMux := runtime.NewServeMux()

	// Act
	gateway.SetHandler(newMux)

	// Assert
	assert.Equal(t, newMux, gateway.mux)
}

func TestGateway_SetHandler_WithHttpHandler(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	originalMux := gateway.mux
	customHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Act
	gateway.SetHandler(customHandler)

	// Assert
	// Когда handler не является *runtime.ServeMux, mux устанавливается в nil
	assert.Nil(t, gateway.mux)
	assert.NotEqual(t, originalMux, gateway.mux)
}

func TestGateway_SetHandler_UpdatesServerHandler(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	gateway.server = &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 30 * time.Second,
	}
	customHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Act
	gateway.SetHandler(customHandler)

	// Assert
	assert.NotNil(t, gateway.server.Handler)
	// Не можем напрямую сравнивать функции, но проверяем что handler установлен
}

func TestGateway_Stop_WithoutServer(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	ctx := context.Background()

	// Act
	err := gateway.Stop(ctx)

	// Assert
	assert.NoError(t, err)
}

func TestGateway_Stop_WithServer(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	gateway.server = &http.Server{
		Addr:              ":0", // Используем случайный порт
		Handler:           http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		ReadHeaderTimeout: 30 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		_ = gateway.server.ListenAndServe()
	}()

	// Даем серверу время на запуск
	time.Sleep(100 * time.Millisecond) // nolint:forbidigo // Необходимо для тестирования запуска сервера

	ctx := context.Background()

	// Act
	err := gateway.Stop(ctx)

	// Assert
	assert.NoError(t, err)
}

func TestGateway_Stop_WithContext(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	gateway.server = &http.Server{
		Addr:              ":0",
		Handler:           http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		ReadHeaderTimeout: 30 * time.Second,
	}

	// Запускаем сервер
	go func() {
		_ = gateway.server.ListenAndServe()
	}()
	time.Sleep(50 * time.Millisecond) // nolint:forbidigo // Необходимо для тестирования

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Act
	err := gateway.Stop(ctx)

	// Assert
	assert.NoError(t, err)
}

func TestGateway_MuxHandlesRequests(t *testing.T) {
	// Arrange
	gateway := NewGateway()

	// Создаем тестовый HTTP запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Act
	gateway.mux.ServeHTTP(w, req)

	// Assert
	// Mux должен вернуть 404 для незарегистрированного маршрута
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGateway_SetHandler_PreservesNilMux(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Act
	gateway.SetHandler(handler)

	// Assert
	assert.Nil(t, gateway.mux, "mux должен быть nil при установке обычного handler")
}

func TestGateway_MultipleSetHandler(t *testing.T) {
	// Arrange
	gateway := NewGateway()
	mux1 := runtime.NewServeMux()
	mux2 := runtime.NewServeMux()

	// Act
	gateway.SetHandler(mux1)
	firstMux := gateway.mux

	gateway.SetHandler(mux2)
	secondMux := gateway.mux

	// Assert
	assert.Equal(t, mux1, firstMux)
	assert.Equal(t, mux2, secondMux)
	assert.NotEqual(t, firstMux, secondMux)
}
