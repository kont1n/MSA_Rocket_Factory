package logger

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SimpleOTLPCore реализует zapcore.Core для отправки логов в OpenTelemetry коллектор.
// Преобразует zap Entry в OpenTelemetry Record и отправляет через OTLP логгер.
type SimpleOTLPCore struct {
	otlpLogger log.Logger
	level      zap.AtomicLevel
}

// NewSimpleOTLPCore создает новый SimpleOTLPCore.
func NewSimpleOTLPCore(otlpLogger log.Logger, level zap.AtomicLevel) *SimpleOTLPCore {
	return &SimpleOTLPCore{
		otlpLogger: otlpLogger,
		level:      level,
	}
}

// Enabled проверяет, включен ли указанный уровень логирования.
func (c *SimpleOTLPCore) Enabled(level zapcore.Level) bool {
	return level >= c.level.Level()
}

// With добавляет поля к core (не используется в OTLP).
func (c *SimpleOTLPCore) With(fields []zap.Field) zapcore.Core {
	return c
}

// Check определяет, нужно ли логировать запись на указанном уровне.
func (c *SimpleOTLPCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}
	return checkedEntry
}

// Write преобразует zap Entry в OpenTelemetry Record и отправляет.
func (c *SimpleOTLPCore) Write(entry zapcore.Entry, fields []zap.Field) error {
	if !c.Enabled(entry.Level) {
		return nil
	}

	// Преобразуем zap поля в OpenTelemetry атрибуты
	attrs := make([]log.KeyValue, 0, len(fields)+3)

	// Добавляем базовые поля
	attrs = append(attrs, log.String("level", entry.Level.String()))
	attrs = append(attrs, log.String("logger", entry.LoggerName))
	attrs = append(attrs, log.String("caller", entry.Caller.String()))

	// Добавляем пользовательские поля
	for _, field := range fields {
		attr := convertZapFieldToOTLP(field)
		if attr.Key != "" {
			attrs = append(attrs, attr)
		}
	}

	// Создаем Record используя правильный конструктор
	record := log.Record{}
	record.SetTimestamp(entry.Time)
	record.SetSeverity(convertZapLevelToOTLP(entry.Level))
	record.SetBody(log.StringValue(entry.Message))
	record.AddAttributes(attrs...)

	// Отправляем Record с контекстом
	c.otlpLogger.Emit(context.Background(), record)

	return nil
}

// Sync принудительно сбрасывает буферы (OTLP SDK делает это автоматически).
func (c *SimpleOTLPCore) Sync() error {
	return nil
}

// convertZapFieldToOTLP преобразует zap.Field в log.KeyValue.
func convertZapFieldToOTLP(field zap.Field) log.KeyValue {
	switch field.Type {
	case zapcore.StringType:
		return log.String(field.Key, field.String)
	case zapcore.Int64Type:
		return log.Int64(field.Key, field.Integer)
	case zapcore.Int32Type:
		return log.Int64(field.Key, int64(field.Integer))
	case zapcore.Uint64Type:
		return log.Int64(field.Key, int64(field.Integer))
	case zapcore.Uint32Type:
		return log.Int64(field.Key, int64(field.Integer))
	case zapcore.Float64Type:
		return log.Float64(field.Key, field.Interface.(float64))
	case zapcore.Float32Type:
		return log.Float64(field.Key, float64(field.Interface.(float32)))
	case zapcore.BoolType:
		return log.Bool(field.Key, field.Integer == 1)
	case zapcore.DurationType:
		return log.String(field.Key, time.Duration(field.Integer).String())
	case zapcore.TimeType:
		return log.String(field.Key, time.Unix(0, field.Integer).Format(time.RFC3339Nano))
	case zapcore.TimeFullType:
		return log.String(field.Key, field.Interface.(time.Time).Format(time.RFC3339Nano))
	default:
		// Для сложных типов используем строковое представление
		return log.String(field.Key, field.String)
	}
}

// convertZapLevelToOTLP преобразует zapcore.Level в log.Severity.
func convertZapLevelToOTLP(level zapcore.Level) log.Severity {
	switch level {
	case zapcore.DebugLevel:
		return log.SeverityDebug
	case zapcore.InfoLevel:
		return log.SeverityInfo
	case zapcore.WarnLevel:
		return log.SeverityWarn
	case zapcore.ErrorLevel:
		return log.SeverityError
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return log.SeverityFatal
	default:
		return log.SeverityInfo
	}
}
