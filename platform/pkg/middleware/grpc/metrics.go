package grpc

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ClientMetrics содержит метрики для gRPC клиентов
type ClientMetrics struct {
	requestsTotal    metric.Int64Counter
	requestDuration  metric.Float64Histogram
	requestErrors    metric.Int64Counter
	connectionStatus metric.Int64UpDownCounter
}

// newClientMetrics создает новый экземпляр метрик для gRPC клиентов
func NewClientMetrics() (*ClientMetrics, error) {
	meter := otel.Meter("grpc-client")

	// Счетчик общего количества запросов
	requestsTotal, err := meter.Int64Counter(
		"grpc_client_requests_total",
		metric.WithDescription("Общее количество gRPC клиентских запросов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Гистограмма времени выполнения запросов
	requestDuration, err := meter.Float64Histogram(
		"grpc_client_request_duration_seconds",
		metric.WithDescription("Время выполнения gRPC клиентских запросов"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик ошибок запросов
	requestErrors, err := meter.Int64Counter(
		"grpc_client_request_errors_total",
		metric.WithDescription("Количество ошибок gRPC клиентских запросов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Gauge статуса соединений
	connectionStatus, err := meter.Int64UpDownCounter(
		"grpc_client_connection_status",
		metric.WithDescription("Статус соединений с gRPC сервисами"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &ClientMetrics{
		requestsTotal:    requestsTotal,
		requestDuration:  requestDuration,
		requestErrors:    requestErrors,
		connectionStatus: connectionStatus,
	}, nil
}

// MetricsInterceptor создает unary client interceptor для мониторинга gRPC запросов
func (m *ClientMetrics) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		startTime := time.Now()

		// Извлекаем имя сервиса из полного пути метода
		service := extractServiceName(method)

		// Выполняем запрос
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Вычисляем время выполнения
		duration := time.Since(startTime)

		// Получаем статус код
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Unknown
			}
		}

		// Формируем атрибуты
		attrs := []attribute.KeyValue{
			attribute.String("service", service),
			attribute.String("method", method),
			attribute.String("status_code", statusCode.String()),
		}

		// Записываем метрики
		m.requestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
		m.requestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

		// Если есть ошибка, записываем её отдельно
		if err != nil {
			attrs = append(attrs, attribute.String("error_type", getErrorType(statusCode)))
			m.requestErrors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}

		return err
	}
}

// RecordConnectionStatus записывает статус соединения с сервисом
func (m *ClientMetrics) RecordConnectionStatus(ctx context.Context, service string, isConnected bool) {
	attrs := []attribute.KeyValue{
		attribute.String("service", service),
	}

	var status int64
	if isConnected {
		status = 1
	}

	m.connectionStatus.Add(ctx, status, metric.WithAttributes(attrs...))
}

// extractServiceName извлекает имя сервиса из полного пути метода
func extractServiceName(method string) string {
	// Пример: /inventory.v1.InventoryService/ListParts -> inventory
	if len(method) > 0 && method[0] == '/' {
		method = method[1:]
	}

	// Находим первую точку
	for i, char := range method {
		if char == '.' {
			return method[:i]
		}
	}

	return "unknown"
}

// getErrorType возвращает тип ошибки на основе gRPC статус кода
func getErrorType(code codes.Code) string {
	switch code {
	case codes.Canceled:
		return "canceled"
	case codes.Unknown:
		return "unknown"
	case codes.InvalidArgument:
		return "invalid_argument"
	case codes.DeadlineExceeded:
		return "deadline_exceeded"
	case codes.NotFound:
		return "not_found"
	case codes.AlreadyExists:
		return "already_exists"
	case codes.PermissionDenied:
		return "permission_denied"
	case codes.ResourceExhausted:
		return "resource_exhausted"
	case codes.FailedPrecondition:
		return "failed_precondition"
	case codes.Aborted:
		return "aborted"
	case codes.OutOfRange:
		return "out_of_range"
	case codes.Unimplemented:
		return "unimplemented"
	case codes.Internal:
		return "internal"
	case codes.Unavailable:
		return "unavailable"
	case codes.DataLoss:
		return "data_loss"
	case codes.Unauthenticated:
		return "unauthenticated"
	default:
		return "unknown"
	}
}
