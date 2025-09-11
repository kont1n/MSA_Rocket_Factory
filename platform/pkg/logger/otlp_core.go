package logger

// OTLP Core Component
//
// Что здесь происходит:
// - record: это одна лог-запись (уровень, сообщение, время, поля-атрибуты).
// - core: «ядро» логгера. Оно решает «принимаю ли я эту запись» и «как её отправлять».
// - tee: «тройник», который раздаёт одну запись сразу нескольким cores.
//
// Интерфейс zapcore.Core (что должен уметь любой core):
// - Enabled(level):решить, писать ли запись этого уровня.
// - With(fields): вернуть копию core с дополнительными полями (мы их учитываем в Write).
// - Check(entry, ce): добавить себя в список получателей записи, если уровень подходит.
// - Write(entry, fields): собрать record и отправить его «куда надо».
// - Sync(): сбросить буферы, если они есть.
//
// Архитектура потока для OTLP:
// zap.Logger -> zapcore.Tee -> SimpleOTLPCore -> OTLP Collector (gRPC)

import (
	"context"
	"time"

	otelLog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap/zapcore"
)

// Таймаут отправки одной записи, чтобы не блокировать приложение
const emitTimeout = 500 * time.Millisecond

// SimpleOTLPCore преобразует zap-записи в OpenTelemetry Records и отправляет их напрямую в OTLP
type SimpleOTLPCore struct {
	otlpLogger otelLog.Logger       // OTLP логгер для отправки записей
	level      zapcore.LevelEnabler // минимальный уровень для записи логов
}

// NewSimpleOTLPCore создает новый OTLP core, работающий напрямую с OTLP-логгером.
func NewSimpleOTLPCore(otlpLogger otelLog.Logger, level zapcore.LevelEnabler) *SimpleOTLPCore {
	return &SimpleOTLPCore{
		otlpLogger: otlpLogger,
		level:      level,
	}
}

// Enabled проверяет, должен ли лог данного уровня быть записан
func (c *SimpleOTLPCore) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

// With создает новый core с дополнительными полями.
// В текущей реализации поля обрабатываются в Write методе,
// поэтому здесь создается копия без изменений.
func (c *SimpleOTLPCore) With(_ []zapcore.Field) zapcore.Core {
	return &SimpleOTLPCore{
		otlpLogger: c.otlpLogger,
		level:      c.level,
	}
}

// Check определяет, должен ли данный лог быть записан данным core.
// Добавляет себя в CheckedEntry если уровень лога соответствует настройкам.
func (c *SimpleOTLPCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

// Write конвертирует zap Entry в OpenTelemetry Record и отправляет в OTLP.
// Пошагово:
//  1. Преобразуем zap-уровень в OTLP Severity (mapZapToOtelSeverity).
//  2. Собираем базовый Record: severity, body=сообщение, timestamp (makeBaseRecord).
//  3. Кодируем zap-поля в OTLP-атрибуты (encodeFieldsToAttrs) и добавляем их в Record.
//  4. Отправляем запись через OTLP-логгер с коротким таймаутом (emitWithTimeout),
//     чтобы не блокировать приложение при сетевых проблемах.
func (c *SimpleOTLPCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	severity := mapZapToOtelSeverity(entry.Level)
	record := makeBaseRecord(entry, severity)
	if len(fields) > 0 {
		attrs := encodeFieldsToAttrs(fields)
		if len(attrs) > 0 {
			record.AddAttributes(attrs...)
		}
	}

	c.emitWithTimeout(record)
	return nil
}

// Sync синхронизация не требуется: батчинг делает OTLP SDK
func (c *SimpleOTLPCore) Sync() error { return nil }

// mapZapToOtelSeverity — отдельная функция преобразования уровня
func mapZapToOtelSeverity(level zapcore.Level) otelLog.Severity {
	switch level {
	case zapcore.DebugLevel:
		return otelLog.SeverityDebug
	case zapcore.InfoLevel:
		return otelLog.SeverityInfo
	case zapcore.WarnLevel:
		return otelLog.SeverityWarn
	case zapcore.ErrorLevel:
		return otelLog.SeverityError
	default:
		return otelLog.SeverityInfo
	}
}

// makeBaseRecord — сборка базового record без атрибутов
func makeBaseRecord(entry zapcore.Entry, sev otelLog.Severity) otelLog.Record {
	r := otelLog.Record{}
	r.SetSeverity(sev)
	r.SetBody(otelLog.StringValue(entry.Message))
	r.SetTimestamp(entry.Time)

	return r
}

// encodeFieldsToAttrs — подготовка атрибутов из zap-полей.
// Используем zapcore.NewMapObjectEncoder(), чтобы безопасно развернуть []zapcore.Field
// в карту ключ→значение. Далее переносим только базовые типы в OTLP KeyValue.
// Неподдерживаемые типы пропускаем (они продолжат жить в stdout части через zap encoder).
func encodeFieldsToAttrs(fields []zapcore.Field) []otelLog.KeyValue {
	if len(fields) == 0 {
		return nil
	}

	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	attrs := make([]otelLog.KeyValue, 0, len(enc.Fields))
	for k, v := range enc.Fields {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, otelLog.String(k, val))
		case bool:
			attrs = append(attrs, otelLog.Bool(k, val))
		case int64:
			attrs = append(attrs, otelLog.Int64(k, val))
		case float64:
			attrs = append(attrs, otelLog.Float64(k, val))
		}
	}

	return attrs
}

// emitWithTimeout — отправка в OTLP с коротким таймаутом
func (c *SimpleOTLPCore) emitWithTimeout(record otelLog.Record) {
	if c.otlpLogger == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), emitTimeout)
	defer cancel()
	c.otlpLogger.Emit(ctx, record)
}
