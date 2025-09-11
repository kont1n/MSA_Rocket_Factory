package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// DefaultCompressor - алгоритм сжатия по умолчанию
	DefaultCompressor = "gzip"
	// DefaultRetryEnabled - включение повторных попыток по умолчанию
	DefaultRetryEnabled = true
	// DefaultRetryInitialInterval - начальный интервал между повторными попытками
	DefaultRetryInitialInterval = 500 * time.Millisecond
	// DefaultRetryMaxInterval - максимальный интервал между повторными попытками
	DefaultRetryMaxInterval = 5 * time.Second
	// DefaultRetryMaxElapsedTime - максимальное время на все повторные попытки
	DefaultRetryMaxElapsedTime = 30 * time.Second
	// DefaultTimeout - таймаут по умолчанию для операций
	DefaultTimeout = 5 * time.Second
)

// serviceName - имя сервиса для трассировки
var serviceName string

type Config interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

// InitTracer инициализирует глобальный трейсер OpenTelemetry.
// Функция возвращает ошибку в случае неудачи инициализации.
func InitTracer(ctx context.Context, cfg Config) error {
	// Сохраняем имя сервиса для использования в спанах
	serviceName = cfg.ServiceName()

	// Создаем экспортер для отправки трейсов в OpenTelemetry Collector через gRPC
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint()), // Адрес коллектора
		otlptracegrpc.WithInsecure(),                        // Отключаем TLS для локальной разработки
		otlptracegrpc.WithTimeout(DefaultTimeout),
		otlptracegrpc.WithCompressor(DefaultCompressor),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         DefaultRetryEnabled,
			InitialInterval: DefaultRetryInitialInterval,
			MaxInterval:     DefaultRetryMaxInterval,
			MaxElapsedTime:  DefaultRetryMaxElapsedTime,
		}),
	)
	if err != nil {
		return err
	}

	// Создаем ресурс с метаданными сервиса
	// Ресурс добавляет атрибуты к каждому трейсу, помогая идентифицировать источник
	attributeResource, err := resource.New(ctx,
		resource.WithAttributes(
			// Используем стандартные атрибуты OpenTelemetry
			semconv.ServiceName(cfg.ServiceName()),
			semconv.ServiceVersion(cfg.ServiceVersion()),
			attribute.String("environment", cfg.Environment()),
		),
		// Автоматически определяем хост, ОС и другие системные атрибуты
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return err
	}

	// Создаем провайдер трейсов с настроенным экспортером и ресурсом
	// BatchSpanProcessor собирает спаны в пакеты для эффективной отправки
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(attributeResource),
		// Настраиваем семплирование трейсов:
		// 1. ParentBased - учитываем решение о семплировании родительского спана
		// 2. TraceIDRatioBased(1.0) - сохраняем 100% трейсов (1.0 = 100%)
		// В продакшене рекомендуется использовать меньший процент (0.1 = 10%)
		// для снижения нагрузки на систему трассировки
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
	)

	// Устанавливаем глобальный провайдер трейсов
	otel.SetTracerProvider(tracerProvider)

	// Настраиваем пропагацию контекста для передачи между сервисами:
	// 1. TraceContext - стандарт W3C для передачи trace ID и parent span ID через HTTP заголовки
	//    Позволяет связать запросы между сервисами в единый трейс
	// 2. Baggage - механизм для передачи дополнительных метаданных между сервисами
	//    Например: user_id, tenant_id, request_id и другие бизнес-контексты
	// Пропагация - это механизм передачи контекста трассировки между сервисами
	// Когда запрос проходит через несколько сервисов, пропагация позволяет:
	// - Сохранить связь между всеми спанами в цепочке вызовов
	// - Передавать дополнительный контекст между сервисами
	// - Обеспечить сквозную трассировку всего запроса
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

// ShutdownTracer закрывает глобальный трейсер OpenTelemetry.
// Функцию следует вызывать при завершении работы приложения.
func ShutdownTracer(ctx context.Context) error {
	provider := otel.GetTracerProvider()
	if provider == nil {
		return nil
	}

	// Приводим к конкретному типу для вызова Shutdown
	tracerProvider, ok := provider.(*sdktrace.TracerProvider)
	if !ok {
		return nil
	}

	// Закрываем провайдер трейсов при выходе
	err := tracerProvider.Shutdown(ctx)
	if err != nil {
		// Ошибки при закрытии не критичны, но могут привести к потере последних трейсов
		return err
	}

	return nil
}

// StartSpan создает новый спан и возвращает его вместе с новым контекстом.
// Это удобная обертка над trace.Tracer.Start, которая использует глобальный трейсер.
//
// Разница между Tracer и TracerProvider:
// 1. TracerProvider - это фабрика трейсеров, которая:
//   - Управляет жизненным циклом трейсеров
//   - Настраивает экспорт трейсов
//   - Контролирует семплирование
//   - Хранит глобальные настройки
//
// 2. Tracer - это конкретный инструмент для создания спанов:
//   - Создает спаны для определенного сервиса/компонента
//   - Управляет связями между спанами
//   - Добавляет атрибуты к спанам
//   - Отслеживает контекст выполнения
//
// Имя трейсера (serviceName):
// - Используется для идентификации источника спанов в системе трассировки
// - Позволяет группировать спаны по сервисам в UI (например, в Jaeger)
// - Если трейсер с таким именем уже существует - он возвращается
// - Если нет - создается новый трейсер с этим именем
// - В нашем случае используется имя сервиса из конфигурации
//
// Создание спана:
// - При создании первого (корневого) спана генерируется новый trace ID
// - Все последующие спаны в цепочке наследуют этот trace ID
// - Trace ID связывает все спаны одного запроса/операции
// - Если в контексте уже есть trace ID, новый не генерируется
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// Получаем трейсер из глобального провайдера
	// Используем имя сервиса из конфигурации для лучшей идентификации в Jaeger
	return otel.Tracer(serviceName).Start(ctx, name, opts...)
}

// SpanFromContext возвращает текущий активный спан из контекста.
// Если спан не существует, возвращается NoopSpan.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext извлекает trace ID из контекста.
// Возвращает строку с ID трейса или пустую строку, если трейс не найден.
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}
