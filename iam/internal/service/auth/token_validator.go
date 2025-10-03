package auth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
)

// validateRefreshTokenWithContext - валидирует refresh токен с контекстом
func (s *JWTService) validateRefreshTokenWithContext(ctx context.Context, tokenString string) (*model.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Строгая проверка алгоритма подписи
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("неожиданный алгоритм подписи: %v", token.Header["alg"])
		}

		// Дополнительная проверка типа метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(s.jwtConfig.RefreshTokenSecret()), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Проверяем, не находится ли токен в blacklist
	if s.blacklistSvc.IsTokenRevoked(ctx, tokenString) {
		return nil, fmt.Errorf("токен отозван")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Проверяем тип токена
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64) // JWT парсит числа как float64
	if !ok {
		return nil, ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userUUID, ok := claims["user_uuid"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	parsedUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &model.Claims{
		UserID:   int64(userID),
		Username: username,
		UserUUID: parsedUUID,
	}, nil
}

// ValidateAccessToken - валидирует access токен (публичный метод)
func (s *JWTService) ValidateAccessToken(tokenString string) (*model.Claims, error) {
	return s.ValidateAccessTokenWithContext(context.Background(), tokenString)
}

// ValidateAccessTokenWithContext - валидирует access токен с контекстом
func (s *JWTService) ValidateAccessTokenWithContext(ctx context.Context, tokenString string) (*model.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Строгая проверка алгоритма подписи
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("неожиданный алгоритм подписи: %v", token.Header["alg"])
		}

		// Дополнительная проверка типа метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(s.jwtConfig.AccessTokenSecret()), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Проверяем, не находится ли токен в blacklist
	if s.blacklistSvc.IsTokenRevoked(ctx, tokenString) {
		return nil, fmt.Errorf("токен отозван")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Проверяем тип токена
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64) // JWT парсит числа как float64
	if !ok {
		return nil, ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userUUID, ok := claims["user_uuid"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	parsedUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &model.Claims{
		UserID:   int64(userID),
		Username: username,
		UserUUID: parsedUUID,
	}, nil
}
