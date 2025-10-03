package v1

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/converter"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	// Создаем спан для операции получения списка деталей
	ctx, span := tracing.StartSpan(ctx, "inventory.list_parts")
	defer tracing.EndSpanWithStatus(span, nil)

	filter := converter.ToModelPart(req)

	// Добавляем атрибуты фильтра к спану
	uuidsCount, namesCount, categoriesCount := 0, 0, 0
	if filter != nil {
		uuidsCount = len(filter.Uuids)
		namesCount = len(filter.Names)
		categoriesCount = len(filter.Categories)
	}
	span.SetAttributes(
		attribute.Int("inventory.filter.uuids_count", uuidsCount),
		attribute.Int("inventory.filter.names_count", namesCount),
		attribute.Int("inventory.filter.categories_count", categoriesCount),
	)

	span.AddEvent("fetching parts from service")
	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Error(ctx, "Failed to get list part",
			zap.Any("filter", filter),
			zap.Error(err),
		)

		if errors.Is(err, model.ErrConvertFromRepo) {
			return nil, status.Errorf(grpcCodes.Internal, "failed to get list parts")
		}

		return nil, status.Error(grpcCodes.Internal, err.Error())
	}

	// Добавляем информацию о результате к спану
	span.SetAttributes(attribute.Int("inventory.result.parts_count", len(*parts)))
	span.AddEvent("parts fetched successfully")

	protoParts := make([]*inventoryV1.Part, 0, len(*parts))
	for _, part := range *parts {
		protoPart := converter.ToProtoPart(&part)
		protoParts = append(protoParts, protoPart)
	}

	span.SetStatus(codes.Ok, "success")
	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
