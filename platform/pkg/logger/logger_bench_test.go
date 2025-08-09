package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func init() {
	// Инициализируем логгер на "мусорный" writer, чтобы не засорять консоль и не тормозить бенчи
	InitForBenchmark()
}

func BenchmarkGlobalLogger(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(ctx, "test message")
	}
}

func BenchmarkWithLogger(b *testing.B) {
	log := With(zap.String("static_field", "static_value"))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(ctx, "test message")
	}
}

func BenchmarkWithContextLogger(b *testing.B) {
	ctx := context.WithValue(context.Background(), traceIDKey, "trace-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithContext(ctx).Info(ctx, "test message")
	}
}

func BenchmarkChainLogger(b *testing.B) {
	ctx := context.WithValue(context.Background(), traceIDKey, "trace-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")

	log := With(zap.String("static_field", "static_value"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(ctx, "test message")
	}
}
