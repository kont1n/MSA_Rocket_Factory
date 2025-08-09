//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Пытаемся загрузить .env файл, но не падаем если его нет
	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "order", ".env"))
	if err != nil {
		logger.Warn(suiteCtx, "Не удалось загрузить .env файл, используем значения по умолчанию", zap.Error(err))
	} else {
		// Устанавливаем переменные в окружение процесса только если файл найден
		for key, value := range envVars {
			_ = os.Setenv(key, value)
		}
	}

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
