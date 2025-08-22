package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

// TokenBlacklistService управляет списком отозванных токенов
type TokenBlacklistService struct {
	cacheRepo repository.SessionCache
}

// NewTokenBlacklistService создает новый сервис blacklist токенов
func NewTokenBlacklistService(cacheRepo repository.SessionCache) *TokenBlacklistService {
	return &TokenBlacklistService{
		cacheRepo: cacheRepo,
	}
}

// RevokeToken добавляет токен в blacklist
func (tbs *TokenBlacklistService) RevokeToken(ctx context.Context, tokenString string) error {
	// Парсим токен для получения информации о времени истечения
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("не удалось разобрать токен: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("неверный формат claims")
	}

	// Получаем время истечения токена
	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("отсутствует время истечения токена")
	}

	expiresAt := time.Unix(int64(exp), 0)
	now := time.Now()

	// Если токен уже истек, не добавляем его в blacklist
	if expiresAt.Before(now) {
		logger.Info(ctx, "Токен уже истек, не добавляем в blacklist")
		return nil
	}

	// Создаем хэш токена для экономии места
	tokenHash := tbs.hashToken(tokenString)

	// TTL = время до истечения токена
	ttl := expiresAt.Sub(now)

	// Создаем фейковую сессию для записи в blacklist
	blacklistEntry := &model.Session{
		UUID:      uuid.MustParse(tokenHash),
		UserUUID:  uuid.New(), // не важно для blacklist
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
	}

	// Добавляем токен в blacklist с TTL
	err = tbs.cacheRepo.Set(ctx, uuid.MustParse(tokenHash), blacklistEntry, ttl)
	if err != nil {
		logger.Error(ctx, "Ошибка добавления токена в blacklist", zap.Error(err))
		return fmt.Errorf("не удалось отозвать токен: %w", err)
	}

	logger.Info(ctx, "✅ Токен добавлен в blacklist", zap.String("token_hash", tokenHash))
	return nil
}

// IsTokenRevoked проверяет, находится ли токен в blacklist
func (tbs *TokenBlacklistService) IsTokenRevoked(ctx context.Context, tokenString string) bool {
	tokenHash := tbs.hashToken(tokenString)

	// Пытаемся получить сессию из кеша по хэшу токена
	session, err := tbs.cacheRepo.GetSessionByUUID(ctx, uuid.MustParse(tokenHash))

	// Если токен найден в blacklist - он отозван
	return err == nil && session != nil
}

// RevokeAllUserTokens отзывает все токены пользователя
func (tbs *TokenBlacklistService) RevokeAllUserTokens(ctx context.Context, userUUID uuid.UUID) error {
	// Создаем UUID из строки для хранения в кеше
	revocationUUID := uuid.New()

	// Создаем фейковую сессию для записи отзыва всех токенов пользователя
	revocationEntry := &model.Session{
		UUID:      revocationUUID,
		UserUUID:  userUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(48 * time.Hour), // время жизни самого долгого refresh токена
	}

	// Устанавливаем маркер на 48 часов (время жизни самого долгого refresh токена)
	err := tbs.cacheRepo.Set(ctx, revocationUUID, revocationEntry, 48*time.Hour)
	if err != nil {
		logger.Error(ctx, "Ошибка отзыва всех токенов пользователя",
			zap.String("user_uuid", userUUID.String()),
			zap.Error(err))
		return fmt.Errorf("не удалось отозвать токены пользователя: %w", err)
	}

	logger.Info(ctx, "✅ Все токены пользователя отозваны", zap.String("user_uuid", userUUID.String()))
	return nil
}

// IsUserTokensRevoked проверяет, отозваны ли все токены пользователя
// Упрощенная версия - для будущих улучшений
func (tbs *TokenBlacklistService) IsUserTokensRevoked(ctx context.Context, userUUID uuid.UUID, tokenIssuedAt time.Time) bool {
	// Пока возвращаем false - можно реализовать позднее, когда будет нужно
	return false
}

// hashToken создает SHA256 хэш токена для использования в качестве ключа
func (tbs *TokenBlacklistService) hashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return fmt.Sprintf("%x", hash)
}
