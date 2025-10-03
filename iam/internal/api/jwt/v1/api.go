package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	jwtV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/jwt/v1"
)

type api struct {
	jwtV1.UnimplementedJWTServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}

func (a *api) Login(ctx context.Context, req *jwtV1.LoginRequest) (*jwtV1.LoginResponse, error) {
	// Валидация запроса
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Выполняем аутентификацию с JWT
	tokenPair, err := a.authService.JWTLogin(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid username or password")
		}
		if errors.Is(err, model.ErrEmptyLogin) || errors.Is(err, model.ErrEmptyPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response := &jwtV1.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	// Добавляем время истечения токенов, если они заданы
	if !tokenPair.AccessTokenExpiresAt.IsZero() {
		response.AccessTokenExpiresAt = timestamppb.New(tokenPair.AccessTokenExpiresAt)
	}
	if !tokenPair.RefreshTokenExpiresAt.IsZero() {
		response.RefreshTokenExpiresAt = timestamppb.New(tokenPair.RefreshTokenExpiresAt)
	}

	return response, nil
}

func (a *api) GetAccessToken(ctx context.Context, req *jwtV1.GetAccessTokenRequest) (*jwtV1.GetAccessTokenResponse, error) {
	// Валидация запроса
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// Получаем новый access токен
	tokenPair, err := a.authService.GetAccessToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	response := &jwtV1.GetAccessTokenResponse{
		AccessToken: tokenPair.AccessToken,
	}

	// Добавляем время истечения токена, если оно задано
	if !tokenPair.AccessTokenExpiresAt.IsZero() {
		response.AccessTokenExpiresAt = timestamppb.New(tokenPair.AccessTokenExpiresAt)
	}

	return response, nil
}

func (a *api) GetRefreshToken(ctx context.Context, req *jwtV1.GetRefreshTokenRequest) (*jwtV1.GetRefreshTokenResponse, error) {
	// Валидация запроса
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// Получаем новый refresh токен
	tokenPair, err := a.authService.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	response := &jwtV1.GetRefreshTokenResponse{
		RefreshToken: tokenPair.RefreshToken,
	}

	// Добавляем время истечения токена, если оно задано
	if !tokenPair.RefreshTokenExpiresAt.IsZero() {
		response.RefreshTokenExpiresAt = timestamppb.New(tokenPair.RefreshTokenExpiresAt)
	}

	return response, nil
}
