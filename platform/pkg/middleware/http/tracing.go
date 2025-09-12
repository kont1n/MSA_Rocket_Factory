package http

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
)

// TracingMiddleware создает HTTP middleware для трассировки входящих запросов.
// Middleware извлекает контекст трассировки из HTTP заголовков и создает новый спан для каждого запроса.
//
// HTTP трейсинг работает следующим образом:
// 1. Клиент отправляет HTTP запрос с заголовками трассировки (traceparent, tracestate)
// 2. Middleware извлекает контекст трассировки из заголовков
// 3. Создается новый спан для HTTP запроса
// 4. Контекст с информацией о трейсе передается дальше по цепочке обработки
// 5. Все последующие операции (gRPC вызовы, работа с БД) наследуют этот контекст
//
// Это обеспечивает сквозную трассировку от HTTP запроса до всех внутренних операций.
func TracingMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Создаем новый спан для HTTP запроса
			// Если в заголовках есть информация о трейсе, спан будет дочерним
			// Если нет - создается корневой спан
			ctx, span := tracing.StartSpan(
				r.Context(),
				r.Method+" "+r.URL.Path,
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			// Добавляем атрибуты HTTP запроса к спану
			// Эти атрибуты помогают идентифицировать запрос в системе трассировки
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.scheme", r.URL.Scheme),
				attribute.String("http.host", r.Host),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.request_id", r.Header.Get("X-Request-ID")),
			)

			// Добавляем информацию о клиенте
			if r.RemoteAddr != "" {
				span.SetAttributes(attribute.String("http.client_ip", r.RemoteAddr))
			}

			// Передаем управление следующему обработчику с обогащенным контекстом
			// Контекст содержит информацию о трейсе, которая будет передана
			// во все последующие операции (gRPC вызовы, работа с БД и т.д.)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
