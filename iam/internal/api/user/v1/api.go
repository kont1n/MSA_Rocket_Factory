package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/converter"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

type api struct {
	iamV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}

func (a *api) Register(ctx context.Context, req *iamV1.RegisterRequest) (*iamV1.RegisterResponse, error) {
	// Конвертируем запрос в доменную модель через конвертер
	registrationInfo := converter.ToModelUserRegistrationInfo(req)
	if registrationInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "user registration info is required")
	}

	// Регистрируем пользователя через сервисный слой
	user, err := a.userService.Register(ctx, registrationInfo)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user with this login already exists")
		}
		if errors.Is(err, model.ErrEmptyLogin) || errors.Is(err, model.ErrEmptyEmail) || errors.Is(err, model.ErrEmptyPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrInvalidEmail) {
			return nil, status.Error(codes.InvalidArgument, "invalid email format")
		}
		if errors.Is(err, model.ErrWeakPassword) {
			return nil, status.Error(codes.InvalidArgument, "password is too weak")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &iamV1.RegisterResponse{
		UserUuid: user.UUID.String(),
	}, nil
}

func (a *api) GetUser(ctx context.Context, req *iamV1.GetUserRequest) (*iamV1.GetUserResponse, error) {
	// Парсим UUID пользователя
	userUUID, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user UUID")
	}

	// Получаем пользователя через сервисный слой
	user, err := a.userService.GetUser(ctx, userUUID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &iamV1.GetUserResponse{
		User: converter.ToProtoUser(user),
	}, nil
}
