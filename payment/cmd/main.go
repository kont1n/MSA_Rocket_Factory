package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	customMiddleware "github.com/kont1n/MSA_Rocket_Factory/payment/internal/api/middleware"
	paymentV1API "github.com/kont1n/MSA_Rocket_Factory/payment/internal/api/payment/v1"
	paymentService "github.com/kont1n/MSA_Rocket_Factory/payment/internal/service/payment"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

const (
	grpcPort = 50052
	httpPort = "8081"

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	log.Printf("Payment service starting...")

	// Занимаем порт для gRPC сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Регистрируем сервис
	service := paymentService.NewService()
	api := paymentV1API.NewAPI(service)

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(grpcServer, api)
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер
	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Запускаем HTTP сервер с gRPC Gateway
	var httpServer *http.Server
	go func() {
		// Создаем контекст с отменой
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем gRPC-gateway мультиплексор
		gwmux := runtime.NewServeMux()

		// Настраиваем опции для соединения с gRPC сервером
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Регистрируем gRPC-gateway хендлеры
		err = paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
			ctx,
			gwmux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("failed to register payment service handler: %v", err)
			return
		}

		// Создаем HTTP роутер
		r := chi.NewRouter()

		// Добавляем middleware
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(10 * time.Second))
		r.Use(customMiddleware.RequestLogger)
		r.Use(cors) // Добавляем CORS middleware

		// Монтируем gRPC-gateway на /v1/payment
		r.Mount("/v1/payment", gwmux)

		// Добавляем health check endpoint
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				log.Printf("failed to write health response: %v", err)
			}
		})

		// Создаем HTTP-сервер
		httpServer = &http.Server{
			Addr:              net.JoinHostPort("localhost", httpPort),
			Handler:           r,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		log.Printf("🚀 HTTP server listening on %s\n", httpPort)
		err = httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ HTTP server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down servers...")

	// Останавливаем HTTP сервер
	if httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err = httpServer.Shutdown(ctx)
		if err != nil {
			log.Printf("❌ Error shutting down HTTP server: %v\n", err)
		} else {
			log.Println("✅ HTTP server stopped")
		}
	}

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()
	log.Println("✅ gRPC server stopped")
}

// CORS middleware для разрешения кросс-доменных запросов
func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
