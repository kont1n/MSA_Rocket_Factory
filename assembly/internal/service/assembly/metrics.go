package assembly

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// assemblyMetrics содержит бизнес метрики для сборки ракет
type assemblyMetrics struct {
	assemblyDuration   metric.Float64Histogram
	rocketsAssembled   metric.Int64Counter
	assemblyErrors     metric.Int64Counter
	assemblyInProgress metric.Int64UpDownCounter
}

// newAssemblyMetrics создает новый экземпляр метрик для сборки
func newAssemblyMetrics() (*assemblyMetrics, error) {
	meter := otel.Meter("assembly-service")

	// Гистограмма длительности сборки
	assemblyDuration, err := meter.Float64Histogram(
		"assembly_duration_seconds",
		metric.WithDescription("Длительность сборки ракет"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.0, 5.0, 10.0, 15.0, 20.0, 25.0, 30.0),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик собранных ракет
	rocketsAssembled, err := meter.Int64Counter(
		"rockets_assembled_total",
		metric.WithDescription("Количество собранных ракет"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик ошибок сборки
	assemblyErrors, err := meter.Int64Counter(
		"assembly_errors_total",
		metric.WithDescription("Количество ошибок при сборке"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик ракет в процессе сборки
	assemblyInProgress, err := meter.Int64UpDownCounter(
		"assembly_in_progress",
		metric.WithDescription("Количество ракет в процессе сборки"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &assemblyMetrics{
		assemblyDuration:   assemblyDuration,
		rocketsAssembled:   rocketsAssembled,
		assemblyErrors:     assemblyErrors,
		assemblyInProgress: assemblyInProgress,
	}, nil
}

// recordAssemblyStart записывает метрики при начале сборки
func (m *assemblyMetrics) recordAssemblyStart(ctx context.Context) {
	attrs := []attribute.KeyValue{
		attribute.String("status", "started"),
	}

	// Увеличиваем счетчик ракет в процессе сборки
	m.assemblyInProgress.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// recordAssemblyComplete записывает метрики при завершении сборки
func (m *assemblyMetrics) recordAssemblyComplete(ctx context.Context, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("status", "completed"),
	}

	// Записываем длительность сборки
	m.assemblyDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	// Увеличиваем счетчик собранных ракет
	m.rocketsAssembled.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Уменьшаем счетчик ракет в процессе сборки
	m.assemblyInProgress.Add(ctx, -1, metric.WithAttributes(attrs...))
}

// recordAssemblyError записывает метрики при ошибке сборки
func (m *assemblyMetrics) recordAssemblyError(ctx context.Context, errorType string) {
	attrs := []attribute.KeyValue{
		attribute.String("error_type", errorType),
		attribute.String("status", "error"),
	}

	// Увеличиваем счетчик ошибок
	m.assemblyErrors.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Уменьшаем счетчик ракет в процессе сборки
	m.assemblyInProgress.Add(ctx, -1, metric.WithAttributes(attrs...))
}
