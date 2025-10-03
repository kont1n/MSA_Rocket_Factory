package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/app"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func init() {
	// В Docker контейнере используем переменные окружения
	// В локальной разработке пытаемся загрузить .env файл
	configPath := "../deploy/compose/assembly/.env"
	if err := config.Load(configPath); err != nil {
		// Если .env файл не найден, пробуем загрузить конфигурацию из переменных окружения
		if err := config.Load(); err != nil {
			panic(fmt.Errorf("failed to load config: %w", err))
		}
	}
}

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
