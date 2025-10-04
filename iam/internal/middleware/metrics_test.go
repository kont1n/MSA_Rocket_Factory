package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func init() {
	// Инициализация логгера для всех тестов
	_ = logger.InitSimple("info", false)
}

func TestNewMetricsMiddleware(t *testing.T) {
	// Act
	middleware := NewMetricsMiddleware("test-service")

	// Assert
	assert.NotNil(t, middleware)
	assert.Equal(t, "test-service", middleware.serviceName)
}

func TestMetricsMiddleware_UnaryServerInterceptor_Success(t *testing.T) {
	// Arrange
	middleware := NewMetricsMiddleware("iam-service")
	interceptor := middleware.UnaryServerInterceptor()

	ctx := context.Background()
	req := "test-request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/iam.v1.AuthService/Login",
	}

	handlerCalled := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return "response", nil
	}

	// Act
	resp, err := interceptor(ctx, req, info, handler)

	// Assert
	require.NoError(t, err)
	assert.True(t, handlerCalled, "Handler должен быть вызван")
	assert.Equal(t, "response", resp)
}

func TestMetricsMiddleware_UnaryServerInterceptor_WithError(t *testing.T) {
	// Arrange
	middleware := NewMetricsMiddleware("iam-service")
	interceptor := middleware.UnaryServerInterceptor()

	ctx := context.Background()
	req := "test-request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/iam.v1.UserService/GetUser",
	}

	expectedErr := status.Error(codes.NotFound, "пользователь не найден")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, expectedErr
	}

	// Act
	resp, err := interceptor(ctx, req, info, handler)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
}

func TestMetricsMiddleware_UnaryServerInterceptor_ContextValues(t *testing.T) {
	// Arrange
	middleware := NewMetricsMiddleware("iam-service")
	interceptor := middleware.UnaryServerInterceptor()

	ctx := context.Background()
	req := "test-request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/iam.v1.AuthService/Whoami",
	}

	var capturedCtx context.Context
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		capturedCtx = ctx
		return "response", nil
	}

	// Act
	_, err := interceptor(ctx, req, info, handler)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, capturedCtx)

	// Проверяем, что в контекст были добавлены значения
	startTime := capturedCtx.Value(requestStartKey)
	assert.NotNil(t, startTime)

	method := capturedCtx.Value(methodKey)
	assert.Equal(t, "/iam.v1.AuthService/Whoami", method)
}

func TestMetricsMiddleware_UnaryServerInterceptor_MultipleErrors(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		handlerErr  error
		expectError bool
	}{
		{
			name:        "success",
			method:      "/iam.v1.AuthService/Login",
			handlerErr:  nil,
			expectError: false,
		},
		{
			name:        "not found error",
			method:      "/iam.v1.UserService/GetUser",
			handlerErr:  status.Error(codes.NotFound, "not found"),
			expectError: true,
		},
		{
			name:        "permission denied",
			method:      "/iam.v1.AuthService/Whoami",
			handlerErr:  status.Error(codes.PermissionDenied, "access denied"),
			expectError: true,
		},
		{
			name:        "internal error",
			method:      "/iam.v1.UserService/UpdateUser",
			handlerErr:  errors.New("internal error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			middleware := NewMetricsMiddleware("iam-service")
			interceptor := middleware.UnaryServerInterceptor()

			ctx := context.Background()
			req := "test-request"
			info := &grpc.UnaryServerInfo{
				FullMethod: tt.method,
			}

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				if tt.handlerErr != nil {
					return nil, tt.handlerErr
				}
				return "success", nil
			}

			// Act
			resp, err := interceptor(ctx, req, info, handler)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
