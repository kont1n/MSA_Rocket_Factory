package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
)

// generateTokenPair - генерирует пару токенов
func (s *JWTService) generateTokenPair(user model.User) (*model.TokenPair, error) {
	accessToken, accessExpiresAt, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

// generateAccessToken - генерирует access токен
func (s *JWTService) generateAccessToken(user model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.jwtConfig.AccessTokenTTL())

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"user_uuid": user.UUID.String(),
		"username":  user.Username,
		"login":     user.Login,
		"exp":       expiresAt.Unix(),
		"iat":       time.Now().Unix(),
		"type":      "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.AccessTokenSecret()))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// generateRefreshToken - генерирует refresh токен
func (s *JWTService) generateRefreshToken(user model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.jwtConfig.RefreshTokenTTL())

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"user_uuid": user.UUID.String(),
		"username":  user.Username,
		"login":     user.Login,
		"exp":       expiresAt.Unix(),
		"iat":       time.Now().Unix(),
		"type":      "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
