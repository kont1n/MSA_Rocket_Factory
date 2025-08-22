// Package fixtures предоставляет централизованные ошибки и error scenarios для тестов
package fixtures

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// ErrorScenario описывает сценарий возникновения ошибки
type ErrorScenario struct {
	Name        string
	Error       error
	Description string
	Category    ErrorCategory
}

// ErrorCategory категории ошибок для группировки
type ErrorCategory string

const (
	CategoryValidation     ErrorCategory = "validation"
	CategoryGRPC           ErrorCategory = "grpc"
	CategoryBusiness       ErrorCategory = "business"
	CategoryInfrastructure ErrorCategory = "infrastructure"
	CategoryDependency     ErrorCategory = "dependency"
)

// Предопределённые ошибки для различных сценариев тестирования

// ValidationErrors - ошибки валидации
var ValidationErrors = struct {
	NilOrder           ErrorScenario
	NilUserUUID        ErrorScenario
	EmptyParts         ErrorScenario
	InvalidUUID        ErrorScenario
	EmptyPaymentMethod ErrorScenario
	TooManyParts       ErrorScenario
}{
	NilOrder: ErrorScenario{
		Name:        "nil_order",
		Error:       errors.New("order cannot be nil"),
		Description: "Заказ не может быть nil",
		Category:    CategoryValidation,
	},
	NilUserUUID: ErrorScenario{
		Name:        "nil_user_uuid",
		Error:       errors.New("user UUID cannot be nil"),
		Description: "UUID пользователя не может быть nil",
		Category:    CategoryValidation,
	},
	EmptyParts: ErrorScenario{
		Name:        "empty_parts",
		Error:       model.ErrPartsSpecified,
		Description: "Список деталей не может быть пустым",
		Category:    CategoryValidation,
	},
	InvalidUUID: ErrorScenario{
		Name:        "invalid_uuid",
		Error:       errors.New("invalid UUID format"),
		Description: "Некорректный формат UUID",
		Category:    CategoryValidation,
	},
	EmptyPaymentMethod: ErrorScenario{
		Name:        "empty_payment_method",
		Error:       errors.New("payment method cannot be empty"),
		Description: "Метод оплаты не может быть пустым",
		Category:    CategoryValidation,
	},
	TooManyParts: ErrorScenario{
		Name:        "too_many_parts",
		Error:       errors.New("too many parts in order"),
		Description: "Слишком много деталей в заказе",
		Category:    CategoryValidation,
	},
}

// GRPCErrors - ошибки gRPC
var GRPCErrors = struct {
	Unavailable      ErrorScenario
	DeadlineExceeded ErrorScenario
	Internal         ErrorScenario
	NotFound         ErrorScenario
	PermissionDenied ErrorScenario
	InvalidArgument  ErrorScenario
	Unauthenticated  ErrorScenario
}{
	Unavailable: ErrorScenario{
		Name:        "grpc_unavailable",
		Error:       status.Error(codes.Unavailable, "service unavailable"),
		Description: "Внешний сервис недоступен",
		Category:    CategoryGRPC,
	},
	DeadlineExceeded: ErrorScenario{
		Name:        "grpc_deadline_exceeded",
		Error:       status.Error(codes.DeadlineExceeded, "timeout"),
		Description: "Превышен timeout запроса",
		Category:    CategoryGRPC,
	},
	Internal: ErrorScenario{
		Name:        "grpc_internal",
		Error:       status.Error(codes.Internal, "internal server error"),
		Description: "Внутренняя ошибка сервера",
		Category:    CategoryGRPC,
	},
	NotFound: ErrorScenario{
		Name:        "grpc_not_found",
		Error:       status.Error(codes.NotFound, "resource not found"),
		Description: "Ресурс не найден",
		Category:    CategoryGRPC,
	},
	PermissionDenied: ErrorScenario{
		Name:        "grpc_permission_denied",
		Error:       status.Error(codes.PermissionDenied, "access denied"),
		Description: "Доступ запрещён",
		Category:    CategoryGRPC,
	},
	InvalidArgument: ErrorScenario{
		Name:        "grpc_invalid_argument",
		Error:       status.Error(codes.InvalidArgument, "invalid request"),
		Description: "Некорректные параметры запроса",
		Category:    CategoryGRPC,
	},
	Unauthenticated: ErrorScenario{
		Name:        "grpc_unauthenticated",
		Error:       status.Error(codes.Unauthenticated, "authentication required"),
		Description: "Требуется аутентификация",
		Category:    CategoryGRPC,
	},
}

// BusinessErrors - бизнес-логические ошибки
var BusinessErrors = struct {
	OrderNotFound     ErrorScenario
	OrderAlreadyPaid  ErrorScenario
	OrderCancelled    ErrorScenario
	OrderAssembled    ErrorScenario
	PartsNotFound     ErrorScenario
	InsufficientStock ErrorScenario
	PaymentFailed     ErrorScenario
	InvalidTransition ErrorScenario
}{
	OrderNotFound: ErrorScenario{
		Name:        "order_not_found",
		Error:       model.ErrOrderNotFound,
		Description: "Заказ не найден",
		Category:    CategoryBusiness,
	},
	OrderAlreadyPaid: ErrorScenario{
		Name:        "order_already_paid",
		Error:       model.ErrPaid,
		Description: "Заказ уже оплачен",
		Category:    CategoryBusiness,
	},
	OrderCancelled: ErrorScenario{
		Name:        "order_cancelled",
		Error:       model.ErrCancelled,
		Description: "Заказ отменён",
		Category:    CategoryBusiness,
	},
	OrderAssembled: ErrorScenario{
		Name:        "order_assembled",
		Error:       errors.New("order is assembled"), // Предполагаемая ошибка
		Description: "Заказ уже собран",
		Category:    CategoryBusiness,
	},
	PartsNotFound: ErrorScenario{
		Name:        "parts_not_found",
		Error:       model.ErrPartsListNotFound,
		Description: "Некоторые детали не найдены",
		Category:    CategoryBusiness,
	},
	InsufficientStock: ErrorScenario{
		Name:        "insufficient_stock",
		Error:       errors.New("insufficient stock for parts"),
		Description: "Недостаточно деталей в наличии",
		Category:    CategoryBusiness,
	},
	PaymentFailed: ErrorScenario{
		Name:        "payment_failed",
		Error:       errors.New("payment processing failed"),
		Description: "Ошибка обработки платежа",
		Category:    CategoryBusiness,
	},
	InvalidTransition: ErrorScenario{
		Name:        "invalid_status_transition",
		Error:       errors.New("invalid order status transition"),
		Description: "Недопустимый переход статуса заказа",
		Category:    CategoryBusiness,
	},
}

// InfrastructureErrors - инфраструктурные ошибки
var InfrastructureErrors = struct {
	DatabaseConnection ErrorScenario
	DatabaseConstraint ErrorScenario
	NetworkTimeout     ErrorScenario
	DiskSpace          ErrorScenario
	Memory             ErrorScenario
}{
	DatabaseConnection: ErrorScenario{
		Name:        "database_connection_failed",
		Error:       errors.New("database connection failed"),
		Description: "Ошибка подключения к базе данных",
		Category:    CategoryInfrastructure,
	},
	DatabaseConstraint: ErrorScenario{
		Name:        "database_constraint_violation",
		Error:       errors.New("database constraint violation"),
		Description: "Нарушение ограничений базы данных",
		Category:    CategoryInfrastructure,
	},
	NetworkTimeout: ErrorScenario{
		Name:        "network_timeout",
		Error:       errors.New("network timeout"),
		Description: "Тайм-аут сетевого соединения",
		Category:    CategoryInfrastructure,
	},
	DiskSpace: ErrorScenario{
		Name:        "insufficient_disk_space",
		Error:       errors.New("insufficient disk space"),
		Description: "Недостаточно места на диске",
		Category:    CategoryInfrastructure,
	},
	Memory: ErrorScenario{
		Name:        "out_of_memory",
		Error:       errors.New("out of memory"),
		Description: "Недостаточно памяти",
		Category:    CategoryInfrastructure,
	},
}

// DependencyErrors - ошибки внешних зависимостей
var DependencyErrors = struct {
	InventoryService    ErrorScenario
	PaymentService      ErrorScenario
	NotificationService ErrorScenario
	AuthService         ErrorScenario
}{
	InventoryService: ErrorScenario{
		Name:        "inventory_service_error",
		Error:       status.Error(codes.Unavailable, "inventory service unavailable"),
		Description: "Сервис инвентаря недоступен",
		Category:    CategoryDependency,
	},
	PaymentService: ErrorScenario{
		Name:        "payment_service_error",
		Error:       status.Error(codes.Internal, "payment service error"),
		Description: "Ошибка платёжного сервиса",
		Category:    CategoryDependency,
	},
	NotificationService: ErrorScenario{
		Name:        "notification_service_error",
		Error:       errors.New("notification service unavailable"),
		Description: "Сервис уведомлений недоступен",
		Category:    CategoryDependency,
	},
	AuthService: ErrorScenario{
		Name:        "auth_service_error",
		Error:       status.Error(codes.Unauthenticated, "authentication service error"),
		Description: "Ошибка сервиса аутентификации",
		Category:    CategoryDependency,
	},
}

// ErrorScenariosByCategory возвращает все сценарии ошибок определённой категории
func ErrorScenariosByCategory(category ErrorCategory) []ErrorScenario {
	var scenarios []ErrorScenario

	allScenarios := []ErrorScenario{
		// Validation errors
		ValidationErrors.NilOrder,
		ValidationErrors.NilUserUUID,
		ValidationErrors.EmptyParts,
		ValidationErrors.InvalidUUID,
		ValidationErrors.EmptyPaymentMethod,
		ValidationErrors.TooManyParts,

		// gRPC errors
		GRPCErrors.Unavailable,
		GRPCErrors.DeadlineExceeded,
		GRPCErrors.Internal,
		GRPCErrors.NotFound,
		GRPCErrors.PermissionDenied,
		GRPCErrors.InvalidArgument,
		GRPCErrors.Unauthenticated,

		// Business errors
		BusinessErrors.OrderNotFound,
		BusinessErrors.OrderAlreadyPaid,
		BusinessErrors.OrderCancelled,
		BusinessErrors.OrderAssembled,
		BusinessErrors.PartsNotFound,
		BusinessErrors.InsufficientStock,
		BusinessErrors.PaymentFailed,
		BusinessErrors.InvalidTransition,

		// Infrastructure errors
		InfrastructureErrors.DatabaseConnection,
		InfrastructureErrors.DatabaseConstraint,
		InfrastructureErrors.NetworkTimeout,
		InfrastructureErrors.DiskSpace,
		InfrastructureErrors.Memory,

		// Dependency errors
		DependencyErrors.InventoryService,
		DependencyErrors.PaymentService,
		DependencyErrors.NotificationService,
		DependencyErrors.AuthService,
	}

	for _, scenario := range allScenarios {
		if scenario.Category == category {
			scenarios = append(scenarios, scenario)
		}
	}

	return scenarios
}

// WrapError обогащает ошибку дополнительным контекстом для тестов
func WrapError(err error, context string, args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(context, args...), err)
}

// IsExpectedError проверяет соответствие ожидаемой ошибки
func IsExpectedError(actual, expected error) bool {
	if expected == nil {
		return actual == nil
	}

	if actual == nil {
		return false
	}

	// Проверка gRPC ошибок
	if status.Code(expected) != codes.OK {
		return status.Code(actual) == status.Code(expected)
	}

	// Проверка errors.Is для wrapped errors
	if errors.Is(actual, expected) {
		return true
	}

	// Проверка по содержимому сообщения
	return actual.Error() == expected.Error()
}

// ErrorAssertionHelper помощник для проверки ошибок в тестах
type ErrorAssertionHelper struct {
	expectedError error
	actualError   error
	testName      string
}

// NewErrorAssertionHelper создаёт новый помощник проверки ошибок
func NewErrorAssertionHelper(testName string) *ErrorAssertionHelper {
	return &ErrorAssertionHelper{testName: testName}
}

// ExpectError устанавливает ожидаемую ошибку
func (h *ErrorAssertionHelper) ExpectError(err error) *ErrorAssertionHelper {
	h.expectedError = err
	return h
}

// ActualError устанавливает полученную ошибку
func (h *ErrorAssertionHelper) ActualError(err error) *ErrorAssertionHelper {
	h.actualError = err
	return h
}

// AssertMatch проверяет соответствие ошибок
func (h *ErrorAssertionHelper) AssertMatch() error {
	if !IsExpectedError(h.actualError, h.expectedError) {
		if h.expectedError == nil {
			return fmt.Errorf("test %s: expected no error, got: %w", h.testName, h.actualError)
		}
		return fmt.Errorf("test %s: expected error %w, got: %w", h.testName, h.expectedError, h.actualError)
	}
	return nil
}

// GetDetailedErrorInfo возвращает детальную информацию об ошибке для диагностики
func GetDetailedErrorInfo(err error) map[string]interface{} {
	if err == nil {
		return map[string]interface{}{
			"error":     nil,
			"type":      "no_error",
			"grpc_code": "OK",
		}
	}

	info := map[string]interface{}{
		"error": err.Error(),
		"type":  fmt.Sprintf("%T", err),
	}

	// Информация о gRPC ошибке
	if st, ok := status.FromError(err); ok {
		info["grpc_code"] = st.Code().String()
		info["grpc_message"] = st.Message()
		info["is_grpc"] = true
	} else {
		info["grpc_code"] = "N/A"
		info["is_grpc"] = false
	}

	// Проверка на wrapped error
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		info["wrapped_error"] = unwrapped.Error()
		info["is_wrapped"] = true
	} else {
		info["is_wrapped"] = false
	}

	return info
}
