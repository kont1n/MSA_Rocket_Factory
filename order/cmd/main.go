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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/kont1n/MSA_Rocket_Factory/order/internal/middleware"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	paymentPort = "50051"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	// Создаем хранилище для данных о погоде
	storage := NewOrderStorage()

	// Создаем gRPC соединение
	conn, err := grpc.NewClient(
		net.JoinHostPort("localhost", paymentPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC клиент для обработки запросов к API платежа
	paymentClient := paymentV1.NewPaymentServiceClient(conn)

	// Создаем обработчик API погоды
	orderHandler := NewOrderHandler(storage, paymentClient)

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
	orders map[uuid.UUID]*orderV1.OrderDto
}

// NewOrderStorage создает новое хранилище данных для заказов
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderV1.OrderDto),
	}
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API заказа
type OrderHandler struct {
	storage *OrderStorage
	paymentClient paymentV1.PaymentServiceClient
}

// NewOrderHandler создает новый обработчик запросов к API заказа
func NewOrderHandler(storage *OrderStorage, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage: storage,
		paymentClient: paymentClient,
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
	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderV1.OrderUUID(order.GetOrderUUID()),
		TotalPrice: orderV1.OptTotalPrice{
			Value: orderV1.TotalPrice(order.GetTotalPrice().Value),
			Set:   order.GetTotalPrice().Set,
		},
	}, nil
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
	var payOrder *orderV1.OrderDto

	// Получаем заказ по UUID
	getOrder, _ := h.GetOrderByUUID(ctx, orderV1.GetOrderByUUIDParams{OrderUUID: params.OrderUUID})
	if getOrder == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprint("Не удалось найти заказ с таким UUID: ", params.OrderUUID),
		}, nil
	}

	// Формируем заказ для оплаты
	payOrder = getOrder.(*orderV1.OrderDto)
	payOrder.PaymentMethod = orderV1.OptPaymentMethod{
		Value: orderV1.PaymentMethod(req.PaymentMethod),
		Set:   true,
	}

	// Оплачиваем заказ с помощью gRPC клиента
	response, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid: payOrder.GetOrderUUID().String(),
		UserUuid:  payOrder.GetUserUUID().String(),
		PaymentMethod: paymentV1.PaymentMethod(payOrder.GetPaymentMethod().Value),
	})

	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Внутренняя ошибка сервиса - не удалось оплатить заказ",
		}, nil
	}

	// Обновляем заказ
	payOrder.TransactionUUID = orderV1.OptUUID{
		Value: uuid.MustParse(response.GetTransactionUuid()),
		Set:   true,
	}
	payOrder.Status = orderV1.OrderStatus("PAID")

	// Сохраняем оплату заказа в хранилище
	h.storage.PayOrder(payOrder)

	return &orderV1.PayOrderResponse{
		TransactionUUID: orderV1.TransactionUUID(payOrder.GetTransactionUUID().Value),
	}, nil
}

// PostOrderCancel обрабатывает запрос отмены заказа
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	var cancelOrder *orderV1.OrderDto

	// Получаем заказ по UUID
	getOrder, _ := h.GetOrderByUUID(ctx, orderV1.GetOrderByUUIDParams{OrderUUID: params.OrderUUID})
	if getOrder == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprint("Не удалось найти заказ с таким UUID: ", params.OrderUUID),
		}, nil
	}

	// Формируем отмену заказа
	cancelOrder = getOrder.(*orderV1.OrderDto)

	// Проверяем статус заказа
	if cancelOrder.GetStatus() == orderV1.OrderStatus("PAID") {
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Заказ уже оплачен и не может быть отменён",
		}, nil
	}
	if cancelOrder.GetStatus() == orderV1.OrderStatus("CANCELLED") {
		return &orderV1.CancelOrderNoContent{}, nil
	}

	// Сохраняем отмену заказа в хранилище
	cancelOrder.Status = orderV1.OrderStatus("CANCELLED")
	h.storage.CancelOrder(cancelOrder)

	return &orderV1.CancelOrderNoContent{}, nil
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

// CreateOrder создает заказ в хранилище
func (s *OrderStorage) CreateOrder(userUUID uuid.UUID, partUUIDs []uuid.UUID) *orderV1.OrderDto {
	newOrder := &orderV1.OrderDto{
		OrderUUID: uuid.New(),
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
		Status:    orderV1.OrderStatus("PENDING_PAYMENT"),
		TotalPrice: orderV1.OptFloat32{
			Value: 0,
			Set:   true,
		},
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[newOrder.OrderUUID] = newOrder
	return newOrder
}

// GetOrder получает информацию о заказе из хранилища
func (s *OrderStorage) GetOrder(orderUUID uuid.UUID) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil
	}
	return order
}

// PayOrder сохраняет оплату заказа в хранилище
func (s *OrderStorage) PayOrder(payOrder *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[payOrder.OrderUUID] = payOrder
}

// CancelOrder сохраняет отмену заказа в хранилище
func (s *OrderStorage) CancelOrder(cancelOrder *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.orders[cancelOrder.OrderUUID] = cancelOrder
}
