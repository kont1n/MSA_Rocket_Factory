package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/kont1n/MSA_Rocket_Factory/order/internal/api/middleware"
	orderV1API "github.com/kont1n/MSA_Rocket_Factory/order/internal/api/order/v1"
	invClient "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/inventory/v1"
	payClient "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/payment/v1"
	oredrRepository "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/inmemory"
	oredrService "github.com/kont1n/MSA_Rocket_Factory/order/internal/service/order"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort      = "8080"
	paymentPort   = "50052"
	inventoryPort = "50051"

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	log.Printf("Order service starting...")

	// Создаем gRPC соединение к API платежа
	paymentConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", paymentPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC соединение к API инвентаря
	inventoryConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", inventoryPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC клиент
	paymentGRPC := paymentV1.NewPaymentServiceClient(paymentConn)
	inventoryGRPC := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := payClient.NewClient(paymentGRPC)
	inventoryClient := invClient.NewClient(inventoryGRPC)

	repo := oredrRepository.NewRepository()
	service := oredrService.NewService(repo, inventoryClient, paymentClient)
	api := orderV1API.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак.
	}
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
