package model

import (
	"fmt"
	"regexp"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// User представляет пользователя в системе
type User struct {
	ID       int64 // Добавлено для совместимости с JWT
	UUID     uuid.UUID
	Login    string
	Username string // Добавлено для совместимости с JWT
	Email    string
	// Password поле удалено для предотвращения случайной утечки
	PasswordHash        string
	NotificationMethods []NotificationMethod
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NotificationMethod представляет способ уведомления пользователя
type NotificationMethod struct {
	ProviderName string // telegram, email, push и т.д.
	Target       string // email адрес, telegram chat id и т.д.
}

// UserRegistrationInfo содержит данные для регистрации нового пользователя
type UserRegistrationInfo struct {
	Login               string
	Email               string
	Password            string
	NotificationMethods []NotificationMethod
}

// TokenPair - пара токенов
type TokenPair struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

// Claims - кастомные claims для JWT
type Claims struct {
	UserID   int64     `json:"user_id"`
	Username string    `json:"username"`
	UserUUID uuid.UUID `json:"user_uuid"`
}

// Константы безопасности
const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MinLoginLength    = 3
	MaxLoginLength    = 50
	MaxEmailLength    = 255
)

var (
	// Регулярное выражение для валидации email
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// Регулярное выражение для валидации логина (только буквы, цифры и подчеркивание)
	loginRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// ValidatePassword проверяет соответствие пароля требованиям безопасности
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("пароль должен содержать минимум %d символов", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("пароль не должен превышать %d символов", MaxPasswordLength)
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("пароль должен содержать минимум одну заглавную букву")
	}
	if !hasLower {
		return fmt.Errorf("пароль должен содержать минимум одну строчную букву")
	}
	if !hasNumber {
		return fmt.Errorf("пароль должен содержать минимум одну цифру")
	}
	if !hasSpecial {
		return fmt.Errorf("пароль должен содержать минимум один специальный символ")
	}

	return nil
}

// ValidateEmail проверяет корректность email адреса
func ValidateEmail(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("email не может быть пустым")
	}
	if len(email) > MaxEmailLength {
		return fmt.Errorf("email не должен превышать %d символов", MaxEmailLength)
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("некорректный формат email")
	}
	return nil
}

// ValidateLogin проверяет корректность логина
func ValidateLogin(login string) error {
	if len(login) < MinLoginLength {
		return fmt.Errorf("логин должен содержать минимум %d символов", MinLoginLength)
	}
	if len(login) > MaxLoginLength {
		return fmt.Errorf("логин не должен превышать %d символов", MaxLoginLength)
	}
	if !loginRegex.MatchString(login) {
		return fmt.Errorf("логин может содержать только буквы, цифры и символ подчеркивания")
	}
	return nil
}

// Validate проверяет корректность данных для регистрации
func (uri *UserRegistrationInfo) Validate() error {
	if err := ValidateLogin(uri.Login); err != nil {
		return fmt.Errorf("ошибка валидации логина: %w", err)
	}

	if err := ValidateEmail(uri.Email); err != nil {
		return fmt.Errorf("ошибка валидации email: %w", err)
	}

	if err := ValidatePassword(uri.Password); err != nil {
		return fmt.Errorf("ошибка валидации пароля: %w", err)
	}

	return nil
}
