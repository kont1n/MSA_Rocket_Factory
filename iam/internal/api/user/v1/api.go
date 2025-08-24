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
	iamV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}

func (a *api) Register(ctx context.Context, req *iamV1.RegisterRequest) (*iamV1.RegisterResponse, error) {
	// Валидация запроса
	if req.Info == nil || req.Info.Info == nil {
		return nil, status.Error(codes.InvalidArgument, "user registration info is required")
	}

	userInfo := req.Info.Info
	if userInfo.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if userInfo.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Info.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Конвертируем способы уведомления
	notificationMethods := make([]model.NotificationMethod, len(userInfo.NotificationMethods))
	for i, method := range userInfo.NotificationMethods {
		notificationMethods[i] = model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	// Создаем информацию для регистрации
	registrationInfo := &model.UserRegistrationInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		Password:            req.Info.Password,
		NotificationMethods: notificationMethods,
	}

	// Регистрируем пользователя
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
	// Валидация UUID пользователя
	userUUID, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user UUID")
	}

	// Получаем пользователя
	user, err := a.userService.GetUser(ctx, userUUID)
	if err != nil {
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

	return &iamV1.GetUserResponse{
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
