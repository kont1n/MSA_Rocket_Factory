package handlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	jwtv1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/jwt/v1"
)

// JWTHandler - gRPC хендлер для JWT сервиса
type JWTHandler struct {
	jwtv1.UnimplementedJWTServiceServer
	jwtService *service.JWTService
}

// NewJWTHandler - создает новый JWT хендлер
func NewJWTHandler(jwtService *service.JWTService) *JWTHandler {
	return &JWTHandler{
		jwtService: jwtService,
	}
}

// Login - обработчик логина
func (h *JWTHandler) Login(ctx context.Context, req *jwtv1.LoginRequest) (*jwtv1.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	tokenPair, err := h.jwtService.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &jwtv1.LoginResponse{
		AccessToken:           tokenPair.AccessToken,
		RefreshToken:          tokenPair.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(tokenPair.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(tokenPair.RefreshTokenExpiresAt),
	}, nil
}

// GetAccessToken - обработчик получения access токена
func (h *JWTHandler) GetAccessToken(ctx context.Context, req *jwtv1.GetAccessTokenRequest) (*jwtv1.GetAccessTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	accessToken, expiresAt, err := h.jwtService.GetAccessToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}
		return nil, status.Error(codes.Internal, "failed to get access token")
	}

	return &jwtv1.GetAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(expiresAt),
	}, nil
}

// GetRefreshToken - обработчик получения refresh токена
func (h *JWTHandler) GetRefreshToken(ctx context.Context, req *jwtv1.GetRefreshTokenRequest) (*jwtv1.GetRefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	refreshToken, expiresAt, err := h.jwtService.GetRefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}
		return nil, status.Error(codes.Internal, "failed to get refresh token")
	}

	return &jwtv1.GetRefreshTokenResponse{
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(expiresAt),
	}, nil
}
