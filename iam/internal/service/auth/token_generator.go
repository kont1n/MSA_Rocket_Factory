package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/olezhek28/microservices-course-examples/week_6/jwt/internal/model"
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
	expiresAt := time.Now().Add(accessTokenTTL)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// generateRefreshToken - генерирует refresh токен
func (s *JWTService) generateRefreshToken(user model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(refreshTokenTTL)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"type":     "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
