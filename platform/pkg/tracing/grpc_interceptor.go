package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Константы для заголовков трассировки
const (
	// TraceIDHeader - заголовок для передачи trace ID
	TraceIDHeader = "x-trace-id"
)

// UnaryServerInterceptor создает gRPC unary interceptor для трассировки входящих запросов.
// Interceptor извлекает контекст трассировки из входящего запроса и создает новый спан для каждого запроса.
//
// Трейсер (Tracer) - это компонент, который создает и управляет спанами. Он отвечает за:
// - Создание новых спанов (начальных или дочерних)
// - Установку ID трейса и связей между спанами
// - Сбор информации о выполнении операций
//
// Пропагатор (Propagator) - это компонент, который отвечает за передачу контекста трассировки
// между сервисами. Он извлекает и внедряет данные трассировки в заголовки запросов,
// что позволяет поддерживать непрерывную трассировку через границы сервисов.
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	// Получаем текущий трейсер и пропагатор
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-interceptor
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Извлекаем метаданные из контекста
		// Incoming context - это контекст, который приходит в запросе от клиента.
		// Он может содержать данные трассировки из предыдущих сервисов в цепочке вызовов.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Извлекаем контекст трассировки из метаданных
		// Пропагатор достает информацию о трейсе из метаданных и добавляет в контекст
		ctx = propagator.Extract(ctx, metadataCarrier(md))

		// Создаем новый спан для текущего запроса
		// Спан (Span) - это единица работы в трассировке, представляющая операцию.
		// Спаны образуют иерархическую структуру (дерево), отражая последовательность операций.
		// Если в контексте уже есть родительский спан, новый спан будет его дочерним.
		// Если родительского спана нет (это первый вызов в цепочке), то создается корневой спан.
		ctx, span := tracer.Start(
			ctx,
			info.FullMethod, // Используем полный метод gRPC как имя спана
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Добавляем trace ID в исходящие метаданные для возврата клиенту
		ctx = AddTraceIDToResponse(ctx)

		// Вызываем обработчик с обогащенным контекстом
		resp, err := handler(ctx, req)
		// Если произошла ошибка, записываем её в спан
		if err != nil {
			span.RecordError(err)
		}

		return resp, err
	}
}

// UnaryClientInterceptor создает gRPC unary interceptor для трассировки исходящих запросов.
// Interceptor добавляет контекст трассировки в исходящий запрос.
//
// Outgoing context - это контекст, который отправляется в запросе к другому сервису.
// В него добавляются данные о текущем трейсе, чтобы следующий сервис мог
// продолжить цепочку трассировки.
func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	// Получаем текущий трейсер и пропагатор
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-interceptor
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Определяем имя спана в зависимости от наличия контекста трассировки
		spanName := formatSpanName(ctx, method)

		// Создаем новый спан с подготовленным именем
		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		// Создаем переносчик метаданных для пропагации трассировки
		carrier := metadataCarrier(extractOutgoingMetadata(ctx))

		// Внедряем контекст трассировки в метаданные
		// Пропагатор добавляет информацию о текущем трейсе в метаданные запроса
		propagator.Inject(ctx, carrier)

		// Обновляем метаданные в контексте
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD(carrier))

		// Вызываем следующий обработчик с обогащенным контекстом
		err := invoker(ctx, method, req, reply, cc, opts...)
		// Если произошла ошибка, записываем её в спан
		if err != nil {
			trace.SpanFromContext(ctx).RecordError(err)
		}

		return err
	}
}

// formatSpanName формирует имя спана в зависимости от наличия контекста трассировки
func formatSpanName(ctx context.Context, method string) string {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return "client." + method
	}

	return method
}

// extractOutgoingMetadata извлекает исходящие метаданные из контекста и создает их копию
func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.New(nil)
	}

	return md.Copy()
}

// GetTraceIDFromContext извлекает trace ID из контекста.
// Полезно для логирования и возврата trace ID клиенту.
func GetTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// AddTraceIDToResponse добавляет trace ID в исходящие метаданные gRPC ответа.
// Это позволяет клиенту получить trace ID для последующего поиска в системе трассировки.
func AddTraceIDToResponse(ctx context.Context) context.Context {
	traceID := GetTraceIDFromContext(ctx)
	if traceID == "" {
		return ctx
	}

	// Используем вспомогательную функцию для извлечения метаданных
	md := extractOutgoingMetadata(ctx)

	md.Set(TraceIDHeader, traceID)
	return metadata.NewOutgoingContext(ctx, md)
}
