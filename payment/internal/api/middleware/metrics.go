package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// HTTPMetrics содержит метрики для HTTP запросов
type HTTPMetrics struct {
	requestCount    metric.Int64Counter
	requestDuration metric.Float64Histogram
	errorCount      metric.Int64Counter
}

// NewHTTPMetrics создает новый экземпляр HTTP метрик
func NewHTTPMetrics() (*HTTPMetrics, error) {
	meter := otel.Meter("payment-service-http")

	// Счетчик HTTP запросов
	requestCount, err := meter.Int64Counter(
		"http_server_request_count",
		metric.WithDescription("Количество HTTP запросов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Гистограмма длительности HTTP запросов
	requestDuration, err := meter.Float64Histogram(
		"http_server_request_duration_seconds",
		metric.WithDescription("Длительность HTTP запросов"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик ошибок HTTP запросов
	errorCount, err := meter.Int64Counter(
		"http_server_errors_count",
		metric.WithDescription("Количество ошибок HTTP запросов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &HTTPMetrics{
		requestCount:    requestCount,
		requestDuration: requestDuration,
		errorCount:      errorCount,
	}, nil
}

// MetricsMiddleware создает middleware для сбора HTTP метрик
func MetricsMiddleware(metrics *HTTPMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Создаем response writer для перехвата статуса
			ww := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Выполняем следующий обработчик
			next.ServeHTTP(ww, r)

			// Вычисляем длительность запроса
			duration := time.Since(start)

			// Нормализуем путь для метрик (заменяем UUID на шаблоны)
			normalizedPath := normalizePath(r.URL.Path)

			// Подготавливаем атрибуты для метрик
			attrs := []attribute.KeyValue{
				attribute.String("http_request_method", r.Method),
				attribute.String("http_route", normalizedPath),
				attribute.String("http_response_status_code", strconv.Itoa(ww.statusCode)),
			}

			// Записываем метрики
			metrics.requestCount.Add(r.Context(), 1, metric.WithAttributes(attrs...))
			metrics.requestDuration.Record(r.Context(), duration.Seconds(), metric.WithAttributes(attrs...))

			// Если статус код указывает на ошибку, увеличиваем счетчик ошибок
			if ww.statusCode >= 400 {
				metrics.errorCount.Add(r.Context(), 1, metric.WithAttributes(attrs...))
			}
		})
	}
}

// metricsResponseWriter обертка для http.ResponseWriter для перехвата статуса
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *metricsResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// normalizePath нормализует HTTP путь для метрик, заменяя UUID на шаблоны
func normalizePath(path string) string {
	// Регулярное выражение для UUID (версия 4)
	uuidRegex := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`)

	// Заменяем UUID на шаблон {uuid}
	normalized := uuidRegex.ReplaceAllString(path, "{uuid}")

	// Специальные случаи для известных endpoints
	switch {
	case strings.HasPrefix(normalized, "/api/v1/payments/{uuid}"):
		return "/api/v1/payments/{uuid}"
	case normalized == "/api/v1/payments":
		return "/api/v1/payments"
	case normalized == "/health":
		return "/health"
	default:
		return normalized
	}
}
