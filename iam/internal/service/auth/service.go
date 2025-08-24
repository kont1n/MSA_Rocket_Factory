package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/config/env"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

var _ def.AuthService = (*service)(nil)

const (
	// Длительность сессии
	sessionDuration = 24 * time.Hour
)

type service struct {
	iamRepository repository.IAMRepository
	jwtService    *JWTService
	blacklistSvc  *TokenBlacklistService
}

func NewService(iamRepository repository.IAMRepository, jwtConfig env.JWTConfig) *service {
	// Создаем blacklist service, используя тот же repository для кеша
	blacklistSvc := NewTokenBlacklistService(iamRepository)

	return &service{
		iamRepository: iamRepository,
		jwtService:    NewJWTService(jwtConfig, blacklistSvc),
		blacklistSvc:  blacklistSvc,
	}
}

func (s *service) Login(ctx context.Context, login, password string) (*model.Session, error) {
	// Валидация входных данных
	if login == "" {
		logger.Warn(ctx, "🚫 Попытка входа с пустым логином")
		return nil, model.ErrEmptyLogin
	}
	if password == "" {
		logger.Warn(ctx, "🚫 Попытка входа с пустым паролем", zap.String("login", login))
		return nil, model.ErrEmptyPassword
	}

	// Получаем пользователя по логину
	user, err := s.iamRepository.GetUserByLogin(ctx, login)
	if err != nil {
		logger.Warn(ctx, "🚫 Неудачная попытка входа: пользователь не найден", zap.String("login", login))
		return nil, model.ErrInvalidCredentials
	}

	// Проверяем пароль
	valid, err := s.verifyPassword(password, user.PasswordHash)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка проверки пароля", zap.String("login", login), zap.Error(err))
		return nil, model.ErrPasswordVerification
	}
	if !valid {
		logger.Warn(ctx, "🚫 Неудачная попытка входа: неверный пароль", zap.String("login", login), zap.String("user_uuid", user.UUID.String()))
		return nil, model.ErrInvalidCredentials
	}

	// Создаем новую сессию
	session := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  user.UUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	// Сохраняем сессию в БД
	session, err = s.iamRepository.CreateSession(ctx, session)
	if err != nil {
		return nil, model.ErrFailedToCreateSession
	}

	// Кешируем сессию в Redis
	if err := s.iamRepository.Set(ctx, session.UUID, session, sessionDuration); err != nil {
		// Логируем ошибку кеширования, но не прерываем выполнение
		// Сессия уже создана в основной БД
		logger.Warn(ctx, "Failed to cache session", zap.Error(err))
	}

	// Логируем успешный вход
	logger.Info(ctx, "✅ Успешный вход в систему",
		zap.String("login", login),
		zap.String("user_uuid", user.UUID.String()),
		zap.String("session_uuid", session.UUID.String()))

	return session, nil
}

func (s *service) Whoami(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, *model.User, error) {
	// Реализация Cache Aside стратегии для получения сессии
	var session *model.Session
	var err error

	// Попытка получить сессию из кеша Redis
	session, err = s.iamRepository.GetSessionFromCache(ctx, sessionUUID)
	if err != nil && !errors.Is(err, model.ErrSessionNotFound) {
		// Если ошибка не "не найдено", то это проблема с кешом, продолжаем работу с БД
		session = nil
	}

	// Если сессии нет в кеше или произошла ошибка кеша, читаем из основной БД
	if session == nil {
		session, err = s.iamRepository.GetSessionByUUID(ctx, sessionUUID)
		if err != nil {
			return nil, nil, err
		}

		// Сохраняем сессию в кеш для следующих запросов
		ttl := time.Until(session.ExpiresAt)
		if ttl > 0 {
			if err := s.iamRepository.Set(ctx, session.UUID, session, ttl); err != nil {
				logger.Warn(ctx, "Failed to cache session", zap.Error(err))
			}
		}
	}

	// Проверяем, не истекла ли сессия
	if session.IsExpired() {
		// Удаляем истёкшую сессию из кеша
		if err := s.iamRepository.Delete(ctx, sessionUUID); err != nil {
			logger.Warn(ctx, "Failed to delete expired session from cache", zap.Error(err))
		}
		return nil, nil, model.ErrSessionExpired
	}

	// Получаем пользователя
	user, err := s.iamRepository.GetUserByUUID(ctx, session.UserUUID)
	if err != nil {
		return nil, nil, err
	}

	return session, user, nil
}

func (s *service) Logout(ctx context.Context, sessionUUID uuid.UUID) error {
	// Удаляем сессию из основной БД
	err := s.iamRepository.DeleteSession(ctx, sessionUUID)
	if err != nil {
		return err
	}

	// Удаляем сессию из кеша (инвалидация)
	if err := s.iamRepository.Delete(ctx, sessionUUID); err != nil {
		logger.Warn(ctx, "Failed to delete session from cache", zap.Error(err))
	}

	return nil
}

// verifyPassword проверяет пароль против хэша
func (s *service) verifyPassword(password, encodedHash string) (bool, error) {
	// Парсим хэш
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, model.ErrPasswordVerification
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, model.ErrPasswordVerification
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, model.ErrPasswordVerification
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, model.ErrPasswordVerification
	}

	// Хэшируем предоставленный пароль с теми же параметрами
	keyLen := len(decodedHash)
	if keyLen < 0 || keyLen > 0xFFFFFFFF {
		return false, model.ErrPasswordVerification
	}
	passwordHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(keyLen))

	// Сравниваем хэши
	return subtle.ConstantTimeCompare(decodedHash, passwordHash) == 1, nil
}

// JWTLogin аутентификация пользователя с возвратом JWT токенов
func (s *service) JWTLogin(ctx context.Context, login, password string) (*model.TokenPair, error) {
	// Валидация входных данных
	if login == "" {
		logger.Warn(ctx, "🚫 JWT: Попытка входа с пустым логином")
		return nil, model.ErrEmptyLogin
	}
	if password == "" {
		logger.Warn(ctx, "🚫 JWT: Попытка входа с пустым паролем", zap.String("login", login))
		return nil, model.ErrEmptyPassword
	}

	// Получаем пользователя по логину
	user, err := s.iamRepository.GetUserByLogin(ctx, login)
	if err != nil {
		logger.Warn(ctx, "🚫 JWT: Неудачная попытка входа: пользователь не найден", zap.String("login", login))
		return nil, model.ErrInvalidCredentials
	}

	// Проверяем пароль
	valid, err := s.verifyPassword(password, user.PasswordHash)
	if err != nil {
		logger.Error(ctx, "❌ JWT: Ошибка проверки пароля", zap.String("login", login), zap.Error(err))
		return nil, model.ErrPasswordVerification
	}
	if !valid {
		logger.Warn(ctx, "🚫 JWT: Неудачная попытка входа: неверный пароль", zap.String("login", login), zap.String("user_uuid", user.UUID.String()))
		return nil, model.ErrInvalidCredentials
	}

	// Заполняем поля для совместимости с JWT
	if user.Username == "" {
		user.Username = user.Login // используем login как username, если нет отдельного username
	}

	// Генерируем JWT токены
	tokenPair, err := s.jwtService.generateTokenPair(*user)
	if err != nil {
		logger.Error(ctx, "❌ JWT: Ошибка генерации токенов", zap.String("login", login), zap.Error(err))
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// Логируем успешный JWT вход
	logger.Info(ctx, "✅ JWT: Успешный вход в систему",
		zap.String("login", login),
		zap.String("user_uuid", user.UUID.String()))

	return tokenPair, nil
}

// GetAccessToken получает новый access токен по refresh токену
func (s *service) GetAccessToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	// Валидируем refresh токен
	claims, err := s.jwtService.validateRefreshTokenWithContext(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Получаем пользователя по UUID из claims
	user, err := s.iamRepository.GetUserByUUID(ctx, claims.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Заполняем поля для совместимости с JWT
	if user.Username == "" {
		user.Username = user.Login
	}

	// Генерируем новый access токен
	accessToken, accessExpiresAt, err := s.jwtService.generateAccessToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken, // возвращаем тот же refresh токен
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: time.Time{}, // не обновляем время истечения refresh токена
	}, nil
}

// GetRefreshToken получает новый refresh токен
func (s *service) GetRefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	// Валидируем текущий refresh токен
	claims, err := s.jwtService.validateRefreshTokenWithContext(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Получаем пользователя по UUID из claims
	user, err := s.iamRepository.GetUserByUUID(ctx, claims.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Заполняем поля для совместимости с JWT
	if user.Username == "" {
		user.Username = user.Login
	}

	// Генерируем новую пару токенов
	tokenPair, err := s.jwtService.generateTokenPair(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return tokenPair, nil
}

// RevokeToken отзывает токен, добавляя его в blacklist
func (s *service) RevokeToken(ctx context.Context, tokenString string) error {
	return s.blacklistSvc.RevokeToken(ctx, tokenString)
}

// RevokeAllUserTokens отзывает все токены пользователя
func (s *service) RevokeAllUserTokens(ctx context.Context, userUUID uuid.UUID) error {
	return s.blacklistSvc.RevokeAllUserTokens(ctx, userUUID)
}
