package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func main() {
	log.Printf("Inventory service starting...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	service := &inventoryService{}
	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	log.Printf("Add Test Data for inventory service")
	TestData(service)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

// GetPart получает информацию о детали по UUID
func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.PartUuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part not found")
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

// ListParts получает список деталей по фильтру
func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	partsFiltered := make([]*inventoryV1.Part, 0)

	s.mu.RLock()
	for _, part := range s.parts {
		partsFiltered = append(partsFiltered, part)
	}
	s.mu.RUnlock()

	filter := req.GetFilter()
	if filter != nil {
		partsFiltered = filtration(filter, partsFiltered)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: partsFiltered,
	}, nil
}

func filtration(filter *inventoryV1.PartsFilter, parts []*inventoryV1.Part) (result []*inventoryV1.Part) {
	log.Printf("filter: %v", filter)
	log.Printf("parts: %v", parts)

	// Создаем мап для фильтрации
	uuidSet := make(map[string]bool)
	for _, uuid := range filter.GetPartUuid() {
		uuidSet[uuid] = true
	}

	nameSet := make(map[string]bool)
	for _, name := range filter.GetPartName() {
		nameSet[name] = true
	}

	categorySet := make(map[inventoryV1.Category]bool)
	for _, category := range filter.GetCategory() {
		categorySet[category] = true
	}

	manufacturerCountrySet := make(map[string]bool)
	for _, manufacturerCountry := range filter.GetManufacturerCountry() {
		manufacturerCountrySet[manufacturerCountry] = true
	}

	tagSet := make(map[string]bool)
	for _, tag := range filter.GetTags() {
		tagSet[tag] = true
	}

	// Фильтруем детали
	log.Printf("uuidSet: %v", uuidSet)
	for _, part := range parts {
		if len(uuidSet) > 0 {
			if _, ok := uuidSet[part.PartUuid]; !ok {
				continue
			}
		}

		if len(nameSet) > 0 {
			if _, ok := nameSet[part.Name]; !ok {
				continue
			}
		}

		if len(categorySet) > 0 {
			if _, ok := categorySet[part.Category]; !ok {
				continue
			}
		}

		if len(manufacturerCountrySet) > 0 {
			if _, ok := manufacturerCountrySet[part.Manufacturer.Country]; !ok {
				continue
			}
		}

		if len(tagSet) > 0 {
			if _, ok := tagSet[part.Tags[0]]; !ok {
				continue
			}
		}

		result = append(result, part)
	}

	return result
}

// TestData Добавление тестовых данных
func TestData(service *inventoryService) {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.parts = map[string]*inventoryV1.Part{
		"d973e963-b7e6-4323-8f4e-4bfd5ab8e834": {
			PartUuid:      "d973e963-b7e6-4323-8f4e-4bfd5ab8e834",
			Name:          "Detail 1",
			Description:   "Detail 1 description",
			Price:         100,
			StockQuantity: 10.0,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 100,
				Width:  100,
				Height: 100,
				Weight: 100,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Country: "China",
				Name:    "Details Fabric",
			},
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
		"d973e963-b7e6-4323-8f4e-4bfd5ab8e835": {
			PartUuid:      "d973e963-b7e6-4323-8f4e-4bfd5ab8e835",
			Name:          "Detail 2",
			Description:   "Detail 2 description",
			Price:         200,
			StockQuantity: 20.0,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 100,
				Width:  100,
				Height: 100,
				Weight: 100,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Country: "USA",
				Name:    "Details Fabric",
			},
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}
}
