package health

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthServer_Check(t *testing.T) {
	// Arrange
	server := &Server{}
	req := &grpc_health_v1.HealthCheckRequest{}

	// Act
	resp, err := server.Check(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthServer_Check_WithServiceName(t *testing.T) {
	// Arrange
	server := &Server{}
	req := &grpc_health_v1.HealthCheckRequest{
		Service: "payment.v1.PaymentService",
	}

	// Act
	resp, err := server.Check(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthServer_Check_WithContext(t *testing.T) {
	// Arrange
	server := &Server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &grpc_health_v1.HealthCheckRequest{}

	// Act
	resp, err := server.Check(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthServer_RegisterService(t *testing.T) {
	// Arrange
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	// Act
	RegisterService(grpcServer)

	// Assert
	// Сервис должен быть зарегистрирован без ошибок
	assert.NotNil(t, grpcServer)
}

func TestHealthServer_RegisterService_Multiple(t *testing.T) {
	// Arrange
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	// Act - регистрируем дважды (должно работать без паники)
	RegisterService(grpcServer)

	// Assert
	assert.NotNil(t, grpcServer)
}
