package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

// LagTracker отслеживает lag для consumer group и записывает метрики
type LagTracker struct {
	client         sarama.Client
	metrics        *KafkaMetrics
	groupID        string
	topics         []string
	updateInterval time.Duration
	stopChan       chan struct{}
	wg             sync.WaitGroup
	logger         interface {
		Info(ctx context.Context, msg string, fields ...zap.Field)
		Error(ctx context.Context, msg string, fields ...zap.Field)
		Debug(ctx context.Context, msg string, fields ...zap.Field)
	}
}

// NewLagTracker создает новый трекер lag'а
func NewLagTracker(client sarama.Client, metrics *KafkaMetrics, groupID string, topics []string, updateInterval time.Duration) *LagTracker {
	return &LagTracker{
		client:         client,
		metrics:        metrics,
		groupID:        groupID,
		topics:         topics,
		updateInterval: updateInterval,
		stopChan:       make(chan struct{}),
		logger:         logger.Logger(),
	}
}

// Start запускает периодическое отслеживание lag'а
func (lt *LagTracker) Start(ctx context.Context) {
	lt.wg.Add(1)
	go lt.trackLag(ctx)
}

// Stop останавливает отслеживание lag'а
func (lt *LagTracker) Stop() {
	close(lt.stopChan)
	lt.wg.Wait()
}

// trackLag периодически отслеживает lag для всех топиков и партиций
func (lt *LagTracker) trackLag(ctx context.Context) {
	defer lt.wg.Done()

	ticker := time.NewTicker(lt.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lt.updateLagMetrics(ctx)
		case <-lt.stopChan:
			return
		}
	}
}

// updateLagMetrics обновляет метрики lag для всех топиков
func (lt *LagTracker) updateLagMetrics(ctx context.Context) {
	// Получаем информацию о consumer group
	coordinator, err := lt.client.Coordinator(lt.groupID)
	if err != nil {
		lt.logger.Error(ctx, "Failed to get coordinator for consumer group",
			zap.String("group_id", lt.groupID),
			zap.Error(err))
		return
	}

	// Получаем информацию о партициях
	partitions := make(map[string][]int32)
	for _, topic := range lt.topics {
		partitions[topic], err = lt.client.Partitions(topic)
		if err != nil {
			lt.logger.Error(ctx, "Failed to get partitions for topic",
				zap.String("topic", topic),
				zap.Error(err))
			continue
		}
	}

	// Для каждой партиции получаем latest offset и committed offset
	for topic, topicPartitions := range partitions {
		for _, partition := range topicPartitions {
			lt.updatePartitionLag(ctx, topic, partition, coordinator)
		}
	}
}

// updatePartitionLag обновляет lag для конкретной партиции
func (lt *LagTracker) updatePartitionLag(ctx context.Context, topic string, partition int32, coordinator *sarama.Broker) {
	// Получаем latest offset (последний доступный offset)
	latestOffset, err := lt.client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		lt.logger.Error(ctx, "Failed to get latest offset",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Error(err))
		return
	}

	// Получаем committed offset (последний обработанный offset)
	request := &sarama.OffsetFetchRequest{
		Version:       1,
		ConsumerGroup: lt.groupID,
	}

	// Добавляем партицию в запрос
	request.AddPartition(topic, partition)

	response, err := coordinator.FetchOffset(request)
	if err != nil {
		lt.logger.Error(ctx, "Failed to fetch committed offset",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Error(err))
		return
	}

	// Вычисляем lag
	var committedOffset int64 = -1
	if block, exists := response.Blocks[topic][partition]; exists && block.Err == sarama.ErrNoError {
		committedOffset = block.Offset
	}

	// Если committed offset не найден, считаем lag равным latest offset
	if committedOffset == -1 {
		committedOffset = 0
	}

	offsetLag := latestOffset - committedOffset
	if offsetLag < 0 {
		offsetLag = 0
	}

	// Получаем timestamp последнего сообщения для расчета lag в секундах
	lagSeconds := lt.calculateLagInSeconds(ctx, topic, partition, latestOffset, committedOffset)

	// Записываем метрику
	lt.metrics.RecordConsumerLag(ctx, topic, partition, lt.groupID, lagSeconds, offsetLag)

	lt.logger.Debug(ctx, "Updated lag metrics",
		zap.String("topic", topic),
		zap.Int32("partition", partition),
		zap.String("group_id", lt.groupID),
		zap.Int64("latest_offset", latestOffset),
		zap.Int64("committed_offset", committedOffset),
		zap.Int64("offset_lag", offsetLag),
		zap.Float64("lag_seconds", lagSeconds))
}

// calculateLagInSeconds рассчитывает lag в секундах
// Это упрощенная реализация - в реальном проекте можно использовать
// timestamp из сообщений для более точного расчета
func (lt *LagTracker) calculateLagInSeconds(_ context.Context, _ string, _ int32, latestOffset, committedOffset int64) float64 {
	offsetLag := latestOffset - committedOffset
	if offsetLag <= 0 {
		return 0
	}

	// Упрощенная оценка lag в секундах
	// Предполагаем, что сообщения приходят со скоростью 100 сообщений в секунду
	// В реальном проекте это значение должно быть настроено на основе мониторинга
	const messagesPerSecond = 100.0
	lagSeconds := float64(offsetLag) / messagesPerSecond

	// Ограничиваем максимальный lag до 3600 секунд (1 час)
	if lagSeconds > 3600 {
		lagSeconds = 3600
	}

	return lagSeconds
}
