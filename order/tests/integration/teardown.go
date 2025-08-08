//go:build integration

package integration

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

// teardownTestEnvironment — освобождает все ресурсы тестового окружения
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "🧹 Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "✅ Тестовое окружение успешно очищено")
}
