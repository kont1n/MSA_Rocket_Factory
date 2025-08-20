package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

type api struct {
	iamV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}

func (a *api) Login(ctx context.Context, req *iamV1.LoginRequest) (*iamV1.LoginResponse, error) {
	// Валидация запроса
	if req.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Выполняем аутентификацию
	session, err := a.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid login or password")
		}
		if errors.Is(err, model.ErrEmptyLogin) || errors.Is(err, model.ErrEmptyPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &iamV1.LoginResponse{
		SessionUuid: session.UUID.String(),
	}, nil
}

func (a *api) Whoami(ctx context.Context, req *iamV1.WhoamiRequest) (*iamV1.WhoamiResponse, error) {
	// Валидация UUID сессии
	sessionUUID, err := uuid.Parse(req.SessionUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session UUID")
	}

	// Получаем информацию о сессии и пользователе
	session, user, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		if errors.Is(err, model.ErrSessionExpired) {
			return nil, status.Error(codes.Unauthenticated, "session expired")
		}
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// Конвертируем способы уведомления
	notificationMethods := make([]*iamV1.NotificationMethod, len(user.NotificationMethods))
	for i, method := range user.NotificationMethods {
		notificationMethods[i] = &iamV1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	return &iamV1.WhoamiResponse{
		Session: &iamV1.Session{
			Uuid:      session.UUID.String(),
			CreatedAt: timestamppb.New(session.CreatedAt),
			UpdatedAt: timestamppb.New(session.UpdatedAt),
			ExpiresAt: timestamppb.New(session.ExpiresAt),
		},
		User: &iamV1.User{
			Uuid: user.UUID.String(),
			Info: &iamV1.UserInfo{
				Login:               user.Login,
				Email:               user.Email,
				NotificationMethods: notificationMethods,
			},
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
