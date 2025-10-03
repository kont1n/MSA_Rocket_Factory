package iam

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

type Client interface {
	RegisterUser(ctx context.Context, login, email, password string, notificationMethods []*iamV1.NotificationMethod) (string, error)
	GetUser(ctx context.Context, userUUID string) (*iamV1.User, error)
	Close() error
}

type client struct {
	conn   *grpc.ClientConn
	client iamV1.UserServiceClient
}

func NewClient(ctx context.Context, target string) (Client, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IAM service: %w", err)
	}

	iamClient := iamV1.NewUserServiceClient(conn)

	return &client{
		conn:   conn,
		client: iamClient,
	}, nil
}

func NewClientWithGRPCClient(grpcClient iamV1.UserServiceClient) Client {
	return &client{
		conn:   nil,
		client: grpcClient,
	}
}

func (c *client) RegisterUser(ctx context.Context, login, email, password string, notificationMethods []*iamV1.NotificationMethod) (string, error) {
	logger.Debug(ctx, "Отправляем gRPC запрос в IAM сервис для регистрации пользователя",
		zap.String("login", login),
		zap.String("email", email))

	req := &iamV1.RegisterRequest{
		Info: &iamV1.UserRegistrationInfo{
			Info: &iamV1.UserInfo{
				Login:               login,
				Email:               email,
				NotificationMethods: notificationMethods,
			},
			Password: password,
		},
	}

	resp, err := c.client.Register(ctx, req)
	if err != nil {
		logger.Error(ctx, "Ошибка при вызове gRPC метода Register в IAM сервисе",
			zap.Error(err),
			zap.String("login", login),
			zap.String("email", email))
		return "", fmt.Errorf("failed to register user via gRPC: %w", err)
	}

	logger.Debug(ctx, "Пользователь успешно зарегистрирован через gRPC вызов в IAM",
		zap.String("login", login),
		zap.String("email", email),
		zap.String("user_uuid", resp.UserUuid))

	return resp.UserUuid, nil
}

func (c *client) GetUser(ctx context.Context, userUUID string) (*iamV1.User, error) {
	logger.Debug(ctx, "Отправляем gRPC запрос в IAM сервис для получения пользователя",
		zap.String("user_uuid", userUUID))

	req := &iamV1.GetUserRequest{
		UserUuid: userUUID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		logger.Error(ctx, "Ошибка при вызове gRPC метода GetUser в IAM сервисе",
			zap.Error(err),
			zap.String("user_uuid", userUUID))
		return nil, fmt.Errorf("failed to get user via gRPC: %w", err)
	}

	logger.Debug(ctx, "Пользователь успешно получен через gRPC вызов в IAM",
		zap.String("user_uuid", userUUID),
		zap.String("login", resp.User.GetInfo().GetLogin()))

	return resp.User, nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
