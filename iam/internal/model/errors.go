package model

import (
	"errors"
	"fmt"
)

// ErrorCode представляет код ошибки для структурированной обработки
type ErrorCode string

const (
	// Коды ошибок пользователей
	ErrCodeUserNotFound     ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserExists       ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeInvalidCreds     ErrorCode = "INVALID_CREDENTIALS" //nolint:gosec
	ErrCodeInvalidUserData  ErrorCode = "INVALID_USER_DATA"
	ErrCodeUserCreateFailed ErrorCode = "USER_CREATE_FAILED"
	ErrCodeUserGetFailed    ErrorCode = "USER_GET_FAILED"
	ErrCodeUserUpdateFailed ErrorCode = "USER_UPDATE_FAILED"

	// Коды ошибок сессий
	ErrCodeSessionNotFound     ErrorCode = "SESSION_NOT_FOUND"
	ErrCodeSessionExpired      ErrorCode = "SESSION_EXPIRED"
	ErrCodeSessionCreateFailed ErrorCode = "SESSION_CREATE_FAILED"
	ErrCodeSessionGetFailed    ErrorCode = "SESSION_GET_FAILED"
	ErrCodeSessionUpdateFailed ErrorCode = "SESSION_UPDATE_FAILED"
	ErrCodeSessionDeleteFailed ErrorCode = "SESSION_DELETE_FAILED"

	// Коды ошибок валидации
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeEmptyLogin       ErrorCode = "EMPTY_LOGIN"
	ErrCodeEmptyPassword    ErrorCode = "EMPTY_PASSWORD"
	ErrCodeEmptyEmail       ErrorCode = "EMPTY_EMAIL"
	ErrCodeInvalidEmail     ErrorCode = "INVALID_EMAIL"
	ErrCodeWeakPassword     ErrorCode = "WEAK_PASSWORD"

	// Коды ошибок безопасности
	ErrCodeHashingFailed      ErrorCode = "HASHING_FAILED"
	ErrCodeVerificationFailed ErrorCode = "VERIFICATION_FAILED"
	ErrCodeTokenInvalid       ErrorCode = "TOKEN_INVALID"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenRevoked       ErrorCode = "TOKEN_REVOKED"

	// Коды системных ошибок
	ErrCodeDatabaseError ErrorCode = "DATABASE_ERROR"
	ErrCodeCacheError    ErrorCode = "CACHE_ERROR"
	ErrCodeInternalError ErrorCode = "INTERNAL_ERROR"
)

// AppError структурированная ошибка приложения
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Cause   error     `json:"-"` // Внутренняя ошибка, не сериализуется
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap позволяет использовать errors.Is и errors.As
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError создает новую структурированную ошибку
func NewAppError(code ErrorCode, message string, details ...string) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// WrapError оборачивает существующую ошибку
func WrapError(code ErrorCode, message string, cause error, details ...string) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// Предопределенные ошибки для совместимости
var (
	// Ошибки пользователей
	ErrUserNotFound       = NewAppError(ErrCodeUserNotFound, "пользователь не найден")
	ErrUserAlreadyExists  = NewAppError(ErrCodeUserExists, "пользователь уже существует")
	ErrInvalidCredentials = NewAppError(ErrCodeInvalidCreds, "неверный логин или пароль")
	ErrInvalidUserData    = NewAppError(ErrCodeInvalidUserData, "некорректные данные пользователя")
	ErrFailedToCreateUser = NewAppError(ErrCodeUserCreateFailed, "не удалось создать пользователя")
	ErrFailedToGetUser    = NewAppError(ErrCodeUserGetFailed, "не удалось получить пользователя")
	ErrFailedToUpdateUser = NewAppError(ErrCodeUserUpdateFailed, "не удалось обновить пользователя")

	// Ошибки сессий
	ErrSessionNotFound       = NewAppError(ErrCodeSessionNotFound, "сессия не найдена")
	ErrSessionExpired        = NewAppError(ErrCodeSessionExpired, "сессия истекла")
	ErrFailedToCreateSession = NewAppError(ErrCodeSessionCreateFailed, "не удалось создать сессию")
	ErrFailedToGetSession    = NewAppError(ErrCodeSessionGetFailed, "не удалось получить сессию")
	ErrFailedToUpdateSession = NewAppError(ErrCodeSessionUpdateFailed, "не удалось обновить сессию")
	ErrFailedToDeleteSession = NewAppError(ErrCodeSessionDeleteFailed, "не удалось удалить сессию")

	// Ошибки валидации
	ErrEmptyLogin    = NewAppError(ErrCodeEmptyLogin, "логин не может быть пустым")
	ErrEmptyPassword = NewAppError(ErrCodeEmptyPassword, "пароль не может быть пустым")
	ErrEmptyEmail    = NewAppError(ErrCodeEmptyEmail, "email не может быть пустым")
	ErrInvalidEmail  = NewAppError(ErrCodeInvalidEmail, "некорректный формат email")
	ErrWeakPassword  = NewAppError(ErrCodeWeakPassword, "пароль слишком слабый")

	// Ошибки хэширования
	ErrFailedToHashPassword = NewAppError(ErrCodeHashingFailed, "не удалось хэшировать пароль")
	ErrPasswordVerification = NewAppError(ErrCodeVerificationFailed, "ошибка проверки пароля")
)

// IsAppError проверяет, является ли ошибка типом AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetErrorCode возвращает код ошибки, если она типа AppError
func GetErrorCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ErrCodeInternalError
}
