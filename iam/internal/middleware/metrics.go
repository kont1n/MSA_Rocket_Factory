package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

// Типы ключей для контекста
type contextKey string

const (
	requestStartKey contextKey = "request_start"
	methodKey       contextKey = "method"
)

// MetricsMiddleware структура для сбора метрик
type MetricsMiddleware struct {
	// В будущем можно добавить Prometheus метрики
	serviceName string
}

// NewMetricsMiddleware создает новый middleware для метрик
func NewMetricsMiddleware(serviceName string) *MetricsMiddleware {
	return &MetricsMiddleware{
		serviceName: serviceName,
	}
}

// UnaryServerInterceptor middleware для unary методов
func (m *MetricsMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Добавляем информацию о запросе в контекст
		ctx = context.WithValue(ctx, requestStartKey, startTime)
		ctx = context.WithValue(ctx, methodKey, info.FullMethod)

		// Выполняем запрос
		resp, err := handler(ctx, req)

		// Вычисляем время выполнения
		duration := time.Since(startTime)

		// Получаем код статуса
		statusCode := status.Code(err)

		// Логируем метрики
		logger.Info(ctx, "📊 Метрики запроса",
			zap.String("service", m.serviceName),
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("status", statusCode.String()),
			zap.Bool("success", err == nil),
		)

		// В будущем здесь можно отправлять метрики в Prometheus
		// m.recordMetrics(info.FullMethod, duration, statusCode)

		return resp, err
	}
}

// В будущем можно добавить интеграцию с Prometheus:
/*
func (m *MetricsMiddleware) recordMetrics(method string, duration time.Duration, code codes.Code) {
	// Prometheus counter для количества запросов
	requestsTotal.WithLabelValues(m.serviceName, method, code.String()).Inc()

	// Prometheus histogram для времени выполнения
	requestDuration.WithLabelValues(m.serviceName, method).Observe(duration.Seconds())
}
*/
