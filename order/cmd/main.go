package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/app"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func main() {
	log.Printf("Order service starting...")

	ctx := context.Background()

	// Загружаем конфигурацию
	err := config.Load("../.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Создаем и запускаем приложение
	orderApp, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	// Запускаем приложение в горутине
	go func() {
		err = orderApp.Run(ctx)
		if err != nil {
			logger.Error(ctx, "failed to run app", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "🛑 Завершение работы сервера...")

	// Закрываем все ресурсы
	err = closer.CloseAll(ctx)
	if err != nil {
		logger.Error(ctx, "failed to close resources", zap.Error(err))
	}

	logger.Info(ctx, "✅ Сервер остановлен")
}
