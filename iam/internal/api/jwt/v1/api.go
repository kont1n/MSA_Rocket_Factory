package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/converter"
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
	// Выполняем аутентификацию через сервисный слой
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

	return converter.ToProtoLoginResponse(tokenPair), nil
}

func (a *api) GetAccessToken(ctx context.Context, req *jwtV1.GetAccessTokenRequest) (*jwtV1.GetAccessTokenResponse, error) {
	// Получаем новый access токен через сервисный слой
	tokenPair, err := a.authService.GetAccessToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	return converter.ToProtoGetAccessTokenResponse(tokenPair), nil
}

func (a *api) GetRefreshToken(ctx context.Context, req *jwtV1.GetRefreshTokenRequest) (*jwtV1.GetRefreshTokenResponse, error) {
	// Получаем новый refresh токен через сервисный слой
	tokenPair, err := a.authService.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	return converter.ToProtoGetRefreshTokenResponse(tokenPair), nil
}
