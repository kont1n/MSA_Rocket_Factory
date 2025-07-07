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
	// Создаем gRPC соединение
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &inventoryService{}
	inventoryV1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
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
	var list []*inventoryV1.Part
	var nameSet []*inventoryV1.Part
	var categorySet []*inventoryV1.Part
	var manufacturerCountrySet []*inventoryV1.Part
	var tagSet []*inventoryV1.Part

	s.mu.RLock()

	// Получаем фильтр
	filter := req.GetFilter()

	// Если фильтр не задан, возвращаем все детали
	if filter == nil {
		for _, part := range s.parts {
			list = append(list, part)
		}
		return &inventoryV1.ListPartsResponse{
			Parts: list,
		}, nil
	}

	// Фильтруем детали по UUID
	uuids := filter.GetPartUuid()
	if len(uuids) > 0 {
		for _, uuid := range uuids {
			part, ok := s.parts[uuid]
			if ok {
				list = append(list, part)
			}
		}
	} else {
		for _, part := range s.parts {
			list = append(list, part)
		}
	}

	s.mu.RUnlock()

	// Фильтруем детали по имени
	names := filter.GetPartName()
	if len(names) > 0 {
		for _, name := range names {
			for _, part := range list {
				if part.GetName() == name {
					nameSet = append(nameSet, part)
				}
			}
		}
	} else {
		nameSet = list
	}

	// Фильтруем детали по категории
	categories := filter.GetCategory()
	if len(categories) > 0 {
		for _, category := range categories {
			for _, part := range list {
				if part.GetCategory() == category {
					categorySet = append(categorySet, part)
				}
			}
		}
	} else {
		categorySet = nameSet
	}

	// Фильтруем детали по стране производителя
	manufacturerCountries := filter.GetManufacturerCountry()
	if len(manufacturerCountries) > 0 {
		for _, manufacturerCountry := range manufacturerCountries {
			for _, part := range list {
				if part.GetManufacturer().GetCountry() == manufacturerCountry {
					manufacturerCountrySet = append(manufacturerCountrySet, part)
				}
			}
		}
	} else {
		manufacturerCountrySet = categorySet
	}

	// Фильтруем детали по тегам
	tags := filter.GetTags()
	if len(tags) > 0 {
		tagSet = tagFilter(manufacturerCountrySet, tags)
	} else {
		tagSet = manufacturerCountrySet
	}

	return &inventoryV1.ListPartsResponse{
		Parts: tagSet,
	}, nil
}

// tagFilter фильтрует детали по тегам
func tagFilter(details []*inventoryV1.Part, tagsFilter []string) (result []*inventoryV1.Part) {
	m := map[string]bool{}
	for _, tag := range tagsFilter {
		m[tag] = true
	}

	for _, detail := range details {
		detailTags := detail.GetTags()
		for _, tag := range detailTags {
			if m[tag] {
				result = append(result, detail)
				break
			}
		}
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
