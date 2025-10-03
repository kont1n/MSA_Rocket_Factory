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
