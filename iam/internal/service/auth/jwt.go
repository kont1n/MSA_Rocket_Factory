package auth

import (
	"errors"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/config/env"
)

var ErrInvalidToken = errors.New("invalid token")

// JWTService - сервис для работы с JWT токенами
type JWTService struct {
	jwtConfig    env.JWTConfig
	blacklistSvc *TokenBlacklistService
}

// NewJWTService - создает новый JWT сервис
func NewJWTService(jwtConfig env.JWTConfig, blacklistSvc *TokenBlacklistService) *JWTService {
	return &JWTService{
		jwtConfig:    jwtConfig,
		blacklistSvc: blacklistSvc,
	}
}
