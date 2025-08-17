package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Запускаем оба Kafka Consumer в горутинах
	go func() {
		err := a.diContainer.OrderPaidConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			logger.Error(ctx, "❌ Ошибка при работе OrderPaid Consumer", zap.Error(err))
		}
	}()

	go func() {
		err := a.diContainer.ShipAssembledConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			logger.Error(ctx, "❌ Ошибка при работе ShipAssembled Consumer", zap.Error(err))
		}
	}()

	// Держим приложение запущенным
	select {
	case <-ctx.Done():
		logger.Info(ctx, "🛑 Получен сигнал завершения работы")
		return nil
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}
