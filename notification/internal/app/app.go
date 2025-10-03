package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
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

	go func() {
		err := a.diContainer.TelegramClient(ctx).Start(ctx)
		if err != nil {
			logger.Error(ctx, "❌ Ошибка при запуске Telegram бота", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "🛑 Получен сигнал завершения работы")
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
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

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

// DiContainer возвращает DI контейнер для тестирования
func (a *App) DiContainer() *diContainer {
	return a.diContainer
}
