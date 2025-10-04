package tracing

import "google.golang.org/grpc/metadata"

// metadataCarrier - это адаптер между gRPC metadata и текстовым отображением OpenTelemetry.
// Реализует интерфейс TextMapCarrier для пропагации контекста трассировки.
type metadataCarrier metadata.MD

// Get возвращает значение для указанного ключа.
func (mc metadataCarrier) Get(key string) string {
	values := metadata.MD(mc).Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// Set устанавливает значение для указанного ключа.
func (mc metadataCarrier) Set(key, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys возвращает список всех ключей.
func (mc metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}

	return keys
}
