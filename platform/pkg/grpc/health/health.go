package health

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server implements the gRPC Health Checking Protocol (GRPC Health v1)
type Server struct {
	grpc_health_v1.UnimplementedHealthServer
}

// Check implements the standard grpc health check protocol
func (s *Server) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch implements the standard grpc health check protocol
func (s *Server) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

// RegisterService registers the health service with the gRPC server
func RegisterService(s *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(s, &Server{})
}
