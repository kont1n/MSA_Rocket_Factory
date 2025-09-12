package kafka

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// KafkaMetrics содержит метрики для Kafka producer и consumer
type KafkaMetrics struct {
	// Producer метрики
	producerMessagesTotal   metric.Int64Counter
	producerMessagesFailed  metric.Int64Counter
	producerMessageDuration metric.Float64Histogram

	// Consumer метрики
	consumerMessagesTotal  metric.Int64Counter
	consumerMessagesFailed metric.Int64Counter
	consumerLagSeconds     metric.Float64Gauge
	consumerOffsetLag      metric.Int64Gauge

	// Consumer Group метрики
	consumerGroupMembers     metric.Int64Gauge
	consumerRebalancingTotal metric.Int64Counter
}

// NewKafkaMetrics создает новый экземпляр метрик для Kafka
func NewKafkaMetrics() (*KafkaMetrics, error) {
	meter := otel.Meter("kafka")

	// Producer метрики
	producerMessagesTotal, err := meter.Int64Counter(
		"kafka_producer_messages_total",
		metric.WithDescription("Общее количество отправленных сообщений"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	producerMessagesFailed, err := meter.Int64Counter(
		"kafka_producer_messages_failed_total",
		metric.WithDescription("Количество неудачных отправок сообщений"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	producerMessageDuration, err := meter.Float64Histogram(
		"kafka_producer_message_duration_seconds",
		metric.WithDescription("Время отправки сообщений"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5),
	)
	if err != nil {
		return nil, err
	}

	// Consumer метрики
	consumerMessagesTotal, err := meter.Int64Counter(
		"kafka_consumer_messages_total",
		metric.WithDescription("Общее количество обработанных сообщений"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	consumerMessagesFailed, err := meter.Int64Counter(
		"kafka_consumer_messages_failed_total",
		metric.WithDescription("Количество ошибок обработки сообщений"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	consumerLagSeconds, err := meter.Float64Gauge(
		"kafka_consumer_lag_seconds",
		metric.WithDescription("Лаг обработки сообщений в секундах"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	consumerOffsetLag, err := meter.Int64Gauge(
		"kafka_consumer_offset_lag",
		metric.WithDescription("Отставание по offset"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Consumer Group метрики
	consumerGroupMembers, err := meter.Int64Gauge(
		"kafka_consumer_group_members",
		metric.WithDescription("Количество участников consumer group"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	consumerRebalancingTotal, err := meter.Int64Counter(
		"kafka_consumer_rebalancing_total",
		metric.WithDescription("Количество ребалансировок consumer group"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &KafkaMetrics{
		producerMessagesTotal:    producerMessagesTotal,
		producerMessagesFailed:   producerMessagesFailed,
		producerMessageDuration:  producerMessageDuration,
		consumerMessagesTotal:    consumerMessagesTotal,
		consumerMessagesFailed:   consumerMessagesFailed,
		consumerLagSeconds:       consumerLagSeconds,
		consumerOffsetLag:        consumerOffsetLag,
		consumerGroupMembers:     consumerGroupMembers,
		consumerRebalancingTotal: consumerRebalancingTotal,
	}, nil
}

// RecordProducerMessage записывает метрики при отправке сообщения
func (m *KafkaMetrics) RecordProducerMessage(ctx context.Context, topic string, partition int32, success bool, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("topic", topic),
		attribute.Int64("partition", int64(partition)),
		attribute.String("status", getProducerStatus(success)),
	}

	// Всегда записываем общее количество сообщений
	m.producerMessagesTotal.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Записываем неудачные отправки отдельно
	if !success {
		m.producerMessagesFailed.Add(ctx, 1, metric.WithAttributes(attrs...))
	}

	m.producerMessageDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// RecordConsumerMessage записывает метрики при обработке сообщения
func (m *KafkaMetrics) RecordConsumerMessage(ctx context.Context, topic string, partition int32, groupID string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("topic", topic),
		attribute.Int64("partition", int64(partition)),
		attribute.String("group_id", groupID),
		attribute.String("status", getConsumerStatus(success)),
	}

	if success {
		m.consumerMessagesTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	} else {
		m.consumerMessagesFailed.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordConsumerLag записывает лаг consumer'а
func (m *KafkaMetrics) RecordConsumerLag(ctx context.Context, topic string, partition int32, groupID string, lagSeconds float64, offsetLag int64) {
	attrs := []attribute.KeyValue{
		attribute.String("topic", topic),
		attribute.Int64("partition", int64(partition)),
		attribute.String("group_id", groupID),
	}

	m.consumerLagSeconds.Record(ctx, lagSeconds, metric.WithAttributes(attrs...))
	m.consumerOffsetLag.Record(ctx, offsetLag, metric.WithAttributes(attrs...))
}

// RecordConsumerGroupMembers записывает количество участников группы
func (m *KafkaMetrics) RecordConsumerGroupMembers(ctx context.Context, groupID string, members int64) {
	attrs := []attribute.KeyValue{
		attribute.String("group_id", groupID),
	}

	m.consumerGroupMembers.Record(ctx, members, metric.WithAttributes(attrs...))
}

// RecordConsumerRebalancing записывает событие ребалансировки
func (m *KafkaMetrics) RecordConsumerRebalancing(ctx context.Context, groupID string) {
	attrs := []attribute.KeyValue{
		attribute.String("group_id", groupID),
	}

	m.consumerRebalancingTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// getProducerStatus возвращает статус producer операции
func getProducerStatus(success bool) string {
	if success {
		return "success"
	}
	return "failed"
}

// getConsumerStatus возвращает статус consumer операции
func getConsumerStatus(success bool) string {
	if success {
		return "success"
	}
	return "failed"
}
