package middleware

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RateLimiter хранит информацию о запросах для каждого IP
type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*clientInfo

	// Настройки
	maxRequests int           // максимальное количество запросов
	window      time.Duration // временное окно
	cleanup     time.Duration // интервал очистки старых записей
}

type clientInfo struct {
	requests []time.Time
	lastSeen time.Time
}

// NewRateLimiter создает новый rate limiter
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:     make(map[string]*clientInfo),
		maxRequests: maxRequests,
		window:      window,
		cleanup:     time.Minute * 5, // очищаем каждые 5 минут
	}

	// Запускаем горутину для периодической очистки
	go rl.cleanupLoop()

	return rl
}

// UnaryServerInterceptor возвращает middleware для unary gRPC методов
func (rl *RateLimiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Получаем IP клиента
		clientIP := rl.getClientIP(ctx)
		if clientIP == "" {
			clientIP = "unknown"
		}

		// Проверяем rate limit только для методов аутентификации
		if rl.isAuthMethod(info.FullMethod) {
			if !rl.allow(clientIP) {
				return nil, status.Error(codes.ResourceExhausted, "слишком много попыток входа, попробуйте позже")
			}
		}

		return handler(ctx, req)
	}
}

// isAuthMethod проверяет, является ли метод методом аутентификации
func (rl *RateLimiter) isAuthMethod(method string) bool {
	authMethods := []string{
		"/iam.v1.AuthService/Login",
		"/jwt.v1.JWTService/Login",
	}

	for _, authMethod := range authMethods {
		if method == authMethod {
			return true
		}
	}

	return false
}

// allow проверяет, можно ли разрешить запрос от данного IP
func (rl *RateLimiter) allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	client, exists := rl.clients[clientIP]
	if !exists {
		client = &clientInfo{
			requests: make([]time.Time, 0),
			lastSeen: now,
		}
		rl.clients[clientIP] = client
	}

	client.lastSeen = now

	// Удаляем старые запросы за пределами окна
	windowStart := now.Add(-rl.window)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range client.requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	client.requests = validRequests

	// Проверяем лимит
	if len(client.requests) >= rl.maxRequests {
		return false
	}

	// Добавляем текущий запрос
	client.requests = append(client.requests, now)

	return true
}

// getClientIP извлекает IP адрес клиента из контекста
func (rl *RateLimiter) getClientIP(ctx context.Context) string {
	// Пытаемся получить IP из metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
			return ips[0]
		}
		if ips := md.Get("x-real-ip"); len(ips) > 0 {
			return ips[0]
		}
	}

	// Если не удалось получить из metadata, возвращаем пустую строку
	// В продакшене можно использовать peer.FromContext для получения адреса
	return ""
}

// cleanupLoop периодически очищает старые записи
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanupOldClients()
	}
}

// cleanupOldClients удаляет клиентов, которые долго не делали запросы
func (rl *RateLimiter) cleanupOldClients() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-time.Hour) // удаляем клиентов, неактивных более часа

	for ip, client := range rl.clients {
		if client.lastSeen.Before(cutoff) {
			delete(rl.clients, ip)
		}
	}
}

// GetStats возвращает текущую статистику rate limiter'a
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	stats := map[string]interface{}{
		"total_clients": len(rl.clients),
		"max_requests":  rl.maxRequests,
		"window":        rl.window.String(),
	}

	return stats
}
