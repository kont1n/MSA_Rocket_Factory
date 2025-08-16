package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/app"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

const configPath = "../deploy/compose/order/.env"

func init() {
	err := config.Load(configPath)
	if err != nil {
		// В контейнере .env файла может не быть, используем переменные окружения
		err = config.Load()
		if err != nil {
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
