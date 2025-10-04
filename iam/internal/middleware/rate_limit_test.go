package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewRateLimiter(t *testing.T) {
	// Act
	rl := NewRateLimiter(10, time.Minute)

	// Assert
	assert.NotNil(t, rl)
	assert.Equal(t, 10, rl.maxRequests)
	assert.Equal(t, time.Minute, rl.window)
	assert.NotNil(t, rl.clients)
}

func TestRateLimiter_IsAuthMethod(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(5, time.Minute)

	tests := []struct {
		name     string
		method   string
		expected bool
	}{
		{
			name:     "auth service login",
			method:   "/iam.v1.AuthService/Login",
			expected: true,
		},
		{
			name:     "jwt service login",
			method:   "/jwt.v1.JWTService/Login",
			expected: true,
		},
		{
			name:     "whoami method",
			method:   "/iam.v1.AuthService/Whoami",
			expected: false,
		},
		{
			name:     "user service",
			method:   "/iam.v1.UserService/GetUser",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := rl.isAuthMethod(tt.method)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(3, time.Second)

	// Act & Assert - первые 3 запроса должны пройти
	assert.True(t, rl.allow("client1"))
	assert.True(t, rl.allow("client1"))
	assert.True(t, rl.allow("client1"))

	// Четвертый запрос должен быть заблокирован
	assert.False(t, rl.allow("client1"))
}

func TestRateLimiter_Allow_DifferentClients(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(2, time.Second)

	// Act & Assert - разные клиенты независимы
	assert.True(t, rl.allow("client1"))
	assert.True(t, rl.allow("client2"))
	assert.True(t, rl.allow("client1"))
	assert.True(t, rl.allow("client2"))

	// Оба клиента достигли лимита
	assert.False(t, rl.allow("client1"))
	assert.False(t, rl.allow("client2"))
}

func TestRateLimiter_Allow_WindowExpires(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(2, 100*time.Millisecond)

	// Act & Assert - заполняем лимит
	assert.True(t, rl.allow("client1"))
	assert.True(t, rl.allow("client1"))
	assert.False(t, rl.allow("client1"))

	// Ждем окончания временного окна
	time.Sleep(150 * time.Millisecond) // nolint:forbidigo // Необходимо для тестирования временных окон

	// Лимит должен обновиться
	assert.True(t, rl.allow("client1"))
}

func TestRateLimiter_UnaryServerInterceptor_NonAuthMethod(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(1, time.Minute)
	interceptor := rl.UnaryServerInterceptor()

	ctx := context.Background()
	req := "test-request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/iam.v1.UserService/GetUser",
	}

	handlerCalled := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return "response", nil
	}

	// Act - для не-auth методов rate limit не применяется
	resp, err := interceptor(ctx, req, info, handler)

	// Assert
	require.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.Equal(t, "response", resp)
}

func TestRateLimiter_UnaryServerInterceptor_AuthMethod_Success(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(5, time.Minute)
	interceptor := rl.UnaryServerInterceptor()

	ctx := context.Background()
	req := "test-request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/iam.v1.AuthService/Login",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	// Act - первый запрос должен пройти
	resp, err := interceptor(ctx, req, info, handler)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "response", resp)
}

func TestRateLimiter_UnaryServerInterceptor_AuthMethod_RateLimited(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(2, time.Minute)
	interceptor := rl.UnaryServerInterceptor()

	ctx := context.Background()
	req := "test-request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/iam.v1.AuthService/Login",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	// Act - первые 2 запроса проходят
	_, err1 := interceptor(ctx, req, info, handler)
	_, err2 := interceptor(ctx, req, info, handler)

	// Третий запрос блокируется
	_, err3 := interceptor(ctx, req, info, handler)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.Error(t, err3)

	st, ok := status.FromError(err3)
	require.True(t, ok)
	assert.Equal(t, codes.ResourceExhausted, st.Code())
	assert.Contains(t, st.Message(), "слишком много попыток входа")
}

func TestRateLimiter_GetClientIP(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(5, time.Minute)

	tests := []struct {
		name     string
		ctx      context.Context // nolint:containedctx // Допустимо в тестовых структурах
		expected string
	}{
		{
			name: "with x-forwarded-for",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"x-forwarded-for", "192.168.1.1",
			)),
			expected: "192.168.1.1",
		},
		{
			name: "with x-real-ip",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				"x-real-ip", "10.0.0.1",
			)),
			expected: "10.0.0.1",
		},
		{
			name:     "without metadata",
			ctx:      context.Background(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := rl.getClientIP(tt.ctx)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRateLimiter_CleanupOldClients(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(5, time.Minute)

	// Добавляем клиентов
	rl.allow("client1")
	rl.allow("client2")

	// Устанавливаем старое время для client1
	rl.mu.Lock()
	rl.clients["client1"].lastSeen = time.Now().Add(-2 * time.Hour)
	rl.mu.Unlock()

	// Act
	rl.cleanupOldClients()

	// Assert
	rl.mu.RLock()
	_, exists1 := rl.clients["client1"]
	_, exists2 := rl.clients["client2"]
	rl.mu.RUnlock()

	assert.False(t, exists1, "Старый клиент должен быть удален")
	assert.True(t, exists2, "Активный клиент должен остаться")
}

func TestRateLimiter_GetStats(t *testing.T) {
	// Arrange
	rl := NewRateLimiter(10, time.Minute)

	rl.allow("client1")
	rl.allow("client2")
	rl.allow("client3")

	// Act
	stats := rl.GetStats()

	// Assert
	assert.Equal(t, 3, stats["total_clients"])
	assert.Equal(t, 10, stats["max_requests"])
	assert.Equal(t, time.Minute.String(), stats["window"])
}
