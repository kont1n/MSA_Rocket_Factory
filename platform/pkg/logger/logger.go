package logger

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

// Константы конфигурации OTLP
const (
	defaultOTLPEndpoint = "localhost:4317" // адрес OTLP коллектора по умолчанию
	shutdownTimeout     = 2 * time.Second  // таймаут для graceful shutdown OTLP provider
)

// Глобальный singleton логгер
var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
	otelProvider *otelLogSdk.LoggerProvider // OTLP provider для graceful shutdown
)

// logger обёртка над zap.Logger с enrich поддержкой контекста
type logger struct {
	zapLogger *zap.Logger
}

// Init инициализирует глобальный логгер с поддержкой OTLP.
// Параметры:
//   - ctx: контекст для операций инициализации
//   - levelStr: уровень логирования ("debug", "info", "warn", "error")
//   - asJSON: формат вывода (true - JSON, false - консольный)
//   - outputs: источники вывода ("stdout", "otlp", "stdout,otlp")
//   - otlpEndpoint: адрес OTLP коллектора (если пустой, используется localhost:4317)
//   - serviceName: имя сервиса для телеметрии
//   - serviceEnvironment: окружение для фильтрации логов
func Init(ctx context.Context, levelStr string, asJSON bool, outputs, otlpEndpoint, serviceName, serviceEnvironment string) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

		// Создаем cores для dual-write архитектуры
		cores := buildCores(ctx, asJSON, outputs, otlpEndpoint, serviceName, serviceEnvironment)

		zapLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(2))

		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})

	return nil
}

// InitSimple инициализирует логгер без OTLP (обратная совместимость)
func InitSimple(levelStr string, asJSON bool) error {
	return Init(context.Background(), levelStr, asJSON, "stdout", "", "", "")
}

// buildCores создает слайс cores для zapcore.Tee на основе outputs.
// Поддерживает "stdout", "otlp", "stdout,otlp".
func buildCores(ctx context.Context, asJSON bool, outputs, otlpEndpoint, serviceName, serviceEnvironment string) []zapcore.Core {
	cores := []zapcore.Core{}

	// Парсим outputs для определения источников вывода
	outputList := parseOutputs(outputs)

	// Добавляем stdout core если указан
	if contains(outputList, "stdout") {
		cores = append(cores, createStdoutCore(asJSON))
	}

	// Добавляем OTLP core если указан
	if contains(outputList, "otlp") {
		if otlpCore := createOTLPCore(ctx, otlpEndpoint, serviceName, serviceEnvironment); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	// Если ничего не указано, используем stdout по умолчанию
	if len(cores) == 0 {
		cores = append(cores, createStdoutCore(asJSON))
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr.
// Поддерживает JSON и консольный формат вывода.
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildProductionEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), dynamicLevel)
}

// createOTLPCore создает core для отправки в OpenTelemetry коллектор.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCore(ctx context.Context, otlpEndpoint, serviceName, serviceEnvironment string) *SimpleOTLPCore {
	if otlpEndpoint == "" {
		otlpEndpoint = defaultOTLPEndpoint
	}
	if serviceName == "" {
		serviceName = "unknown-service"
	}
	if serviceEnvironment == "" {
		serviceEnvironment = "dev"
	}

	otlpLogger, err := createOTLPLogger(ctx, otlpEndpoint, serviceName, serviceEnvironment)
	if err != nil {
		// Логирование ошибки невозможно, так как логгер еще не инициализирован
		return nil
	}

	// Прямо передаём OTLP-логгер в core. Буферизацию делает OTLP SDK (BatchProcessor).
	return NewSimpleOTLPCore(otlpLogger, dynamicLevel)
}

// createOTLPLogger создает OTLP логгер с настроенным экспортером и ресурсами.
// Использует BatchProcessor для эффективной отправки логов.
func createOTLPLogger(ctx context.Context, endpoint, serviceName, serviceEnvironment string) (otelLog.Logger, error) {
	exporter, err := createOTLPExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createResource(ctx, serviceName, serviceEnvironment)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider // сохраняем для shutdown

	return provider.Logger("app"), nil
}

// createOTLPExporter создает gRPC экспортер для OTLP коллектора
func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // для разработки, в продакшене следует использовать TLS
	)
}

// createResource создает метаданные сервиса для телеметрии
func createResource(ctx context.Context, serviceName, serviceEnvironment string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", serviceEnvironment),
		),
	)
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",                 // время
		LevelKey:       "level",                     // уровень логирования
		NameKey:        "logger",                    // имя логгера, если используется
		CallerKey:      "caller",                    // откуда вызван лог
		MessageKey:     "message",                   // текст сообщения
		StacktraceKey:  "stacktrace",                // стектрейс для ошибок
		LineEnding:     zapcore.DefaultLineEnding,   // перенос строки
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // читаемый ISO 8601 формат
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // короткий caller
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// SetLevel динамически меняет уровень логирования
func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}

	dynamicLevel.SetLevel(parseLevel(levelStr))
}

func InitForBenchmark() {
	core := zapcore.NewNopCore()

	globalLogger = &logger{
		zapLogger: zap.New(core),
	}
}

// logger возвращает глобальный enrich-aware логгер
func Logger() *logger {
	return globalLogger
}

// NopLogger устанавливает глобальный логгер в no-op режим.
// Идеально для юнит-тестов.
func SetNopLogger() {
	globalLogger = &logger{
		zapLogger: zap.NewNop(),
	}
}

// Sync сбрасывает буферы логгера
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}

// Close корректно завершает работу логгера.
// Останавливает OTLP provider с таймаутом для отправки оставшихся логов.
func Close() error {
	if otelProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := otelProvider.Shutdown(ctx); err != nil {
			// Логирование ошибки невозможно, так как логгер уже завершается
			_ = err
		}
	}

	return nil
}

// With создает новый enrich-aware логгер с дополнительными полями
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext создает enrich-aware логгер с контекстом
func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// Debug enrich-aware debug log
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}

// Info enrich-aware info log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}

// Warn enrich-aware warn log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}

// Error enrich-aware error log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}

// Fatal enrich-aware fatal log
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}

// Instance methods для enrich loggers (logger)

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// parseLevel конвертирует строковый уровень в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsFromContext вытаскивает enrich-поля из контекста
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}

// parseOutputs парсит строку outputs в слайс источников вывода
func parseOutputs(outputs string) []string {
	if outputs == "" {
		return []string{"stdout"}
	}

	parts := strings.Split(outputs, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// contains проверяет, содержится ли элемент в слайсе
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
