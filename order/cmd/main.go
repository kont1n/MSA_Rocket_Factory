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
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	customMiddleware "github.com/kont1n/MSA_Rocket_Factory/order/internal/middleware"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type Order struct {
	orderV1.GetOrderResponse
}

func main() {
	// Создаем хранилище для данных о погоде
	storage := NewOrderStorage()

	// Создаем обработчик API погоды
	orderHandler := NewOrderHandler(storage)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
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

// OrderStorage представляет потокобезопасное хранилище данных для заказов
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*Order
}

// NewOrderStorage создает новое хранилище данных для заказов
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*Order),
	}
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API заказа
type OrderHandler struct {
	storage *OrderStorage
}

// NewOrderHandler создает новый обработчик запросов к API заказа
func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

// CreateOrder обрабатывает запрос создания заказа
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order := h.storage.CreateOrder(uuid.UUID(req.UserUUID), req.PartUuids)
	if order == nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Внутренняя ошибка сервиса",
		}, nil
	}
	return order, nil
}

// GetOrderByUUID обрабатывает запрос получения информации о заказе по UUID
func (h *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprint("Не удалось найти заказ с таким UUID: ", params.OrderUUID),
		}, nil
	}
	return order, nil
}

// PostOrderPayment обрабатывает запрос оплаты заказа
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	return nil, nil
}

// PostOrderCancel обрабатывает запрос отмены заказа
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	return nil, nil
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	code := orderV1.OptInt{}
	code.SetTo(http.StatusInternalServerError)

	message := orderV1.OptString{}
	message.SetTo(err.Error())

	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    code,
			Message: message,
		},
	}
}

// CreateOrder создает заказ
func (s *OrderStorage) CreateOrder(userUUID uuid.UUID, partUUIDs []uuid.UUID) *Order {

	addOrder := &Order{}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[addOrder.OrderUUID] = addOrder
	return addOrder
}

// GetOrder получает информацию о заказе
func (s *OrderStorage) GetOrder(orderUUID uuid.UUID) *Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil
	}
	return order
}

// PayOrder оплачивает заказ
func (s *OrderStorage) PayOrder(orderUUID string, paymentMethod int) *Order {
	return nil
}

// CancelOrder отменяет заказ
func (s *OrderStorage) CancelOrder(orderUUID string) *Order {
	return nil
}
