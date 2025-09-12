package order

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	grpcMetrics "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/middleware/grpc"
	kafkaMiddleware "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/middleware/kafka"
)

// orderMetrics содержит бизнес метрики для заказов
type orderMetrics struct {
	ordersTotal     metric.Int64Counter
	ordersRevenue   metric.Float64Counter
	ordersCreated   metric.Int64Counter
	ordersPaid      metric.Int64Counter
	ordersCancelled metric.Int64Counter
	orderValue      metric.Float64Histogram
	orderDuration   metric.Float64Histogram
	grpcMetrics     *grpcMetrics.ClientMetrics
	kafkaMetrics    *kafkaMiddleware.KafkaMetrics
}

// newOrderMetrics создает новый экземпляр метрик для заказов
func newOrderMetrics() (*orderMetrics, error) {
	meter := otel.Meter("order-service")

	// Счетчик общего количества заказов
	ordersTotal, err := meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Общее количество заказов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик выручки по валютам
	ordersRevenue, err := meter.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Суммарная выручка от заказов"),
		metric.WithUnit("currency"),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик созданных заказов
	ordersCreated, err := meter.Int64Counter(
		"orders_created_total",
		metric.WithDescription("Количество созданных заказов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик оплаченных заказов
	ordersPaid, err := meter.Int64Counter(
		"orders_paid_total",
		metric.WithDescription("Количество оплаченных заказов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Счетчик отмененных заказов
	ordersCancelled, err := meter.Int64Counter(
		"orders_cancelled_total",
		metric.WithDescription("Количество отмененных заказов"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Гистограмма стоимости заказов
	orderValue, err := meter.Float64Histogram(
		"order_value",
		metric.WithDescription("Стоимость заказов"),
		metric.WithUnit("currency"),
		metric.WithExplicitBucketBoundaries(0, 1000, 5000, 10000, 50000, 100000, 500000, 1000000, 5000000, 10000000, 50000000, 100000000, 500000000, 1000000000),
	)
	if err != nil {
		return nil, err
	}

	// Гистограмма времени обработки заказов
	orderDuration, err := meter.Float64Histogram(
		"order_duration_seconds",
		metric.WithDescription("Время обработки заказов"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0, 0.1, 0.5, 1, 2, 5, 10, 30, 60, 300),
	)
	if err != nil {
		return nil, err
	}

	// Создаем gRPC метрики
	grpcMetrics, err := grpcMetrics.NewClientMetrics()
	if err != nil {
		return nil, err
	}

	// Создаем Kafka метрики
	kafkaMetrics, err := kafkaMiddleware.NewKafkaMetrics()
	if err != nil {
		return nil, err
	}

	return &orderMetrics{
		ordersTotal:     ordersTotal,
		ordersRevenue:   ordersRevenue,
		ordersCreated:   ordersCreated,
		ordersPaid:      ordersPaid,
		ordersCancelled: ordersCancelled,
		orderValue:      orderValue,
		orderDuration:   orderDuration,
		grpcMetrics:     grpcMetrics,
		kafkaMetrics:    kafkaMetrics,
	}, nil
}

// recordOrderCreated записывает метрики при создании заказа
func (m *orderMetrics) recordOrderCreated(ctx context.Context, value float64, currency string) {
	attrs := []attribute.KeyValue{
		attribute.String("currency", currency),
		attribute.String("status", "created"),
	}

	// Увеличиваем счетчики
	m.ordersTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.ordersCreated.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.ordersRevenue.Add(ctx, value, metric.WithAttributes(attrs...))

	// Записываем стоимость в гистограмму
	m.orderValue.Record(ctx, value, metric.WithAttributes(attrs...))
}

// recordOrderPaid записывает метрики при оплате заказа
func (m *orderMetrics) recordOrderPaid(ctx context.Context, value float64, currency string) {
	attrs := []attribute.KeyValue{
		attribute.String("currency", currency),
		attribute.String("status", "paid"),
	}

	// Увеличиваем счетчик оплаченных заказов
	m.ordersPaid.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// recordOrderCancelled записывает метрики при отмене заказа
func (m *orderMetrics) recordOrderCancelled(ctx context.Context, value float64, currency string) {
	attrs := []attribute.KeyValue{
		attribute.String("currency", currency),
		attribute.String("status", "cancelled"),
	}

	// Увеличиваем счетчик отмененных заказов
	m.ordersCancelled.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// recordOrderDuration записывает время обработки заказа
func (m *orderMetrics) recordOrderDuration(ctx context.Context, duration time.Duration, operation string) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
	}

	// Записываем время в гистограмму
	m.orderDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// GetGRPCMetrics возвращает gRPC метрики для использования в middleware
func (m *orderMetrics) GetGRPCMetrics() *grpcMetrics.ClientMetrics {
	return m.grpcMetrics
}

// GetKafkaMetrics возвращает Kafka метрики для использования в producer/consumer
func (m *orderMetrics) GetKafkaMetrics() *kafkaMiddleware.KafkaMetrics {
	return m.kafkaMetrics
}
