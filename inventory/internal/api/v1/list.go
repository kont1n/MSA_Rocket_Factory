package v1

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
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
	defer span.End()

	filter := converter.ToModelPart(req)

	// Добавляем атрибуты к спану
	if filter != nil {
		span.SetAttributes(
			attribute.Int("inventory.filter.uuids_count", len(filter.Uuids)),
			attribute.Int("inventory.filter.names_count", len(filter.Names)),
			attribute.Int("inventory.filter.categories_count", len(filter.Categories)),
		)
	} else {
		span.SetAttributes(
			attribute.Int("inventory.filter.uuids_count", 0),
			attribute.Int("inventory.filter.names_count", 0),
			attribute.Int("inventory.filter.categories_count", 0),
		)
	}

	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "Failed to get list part",
			zap.Any("filter", filter),
			zap.Error(err),
		)

		if errors.Is(err, model.ErrConvertFromRepo) {
			return nil, status.Errorf(codes.Internal, "failed to get list parts")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	// Добавляем информацию о результате к спану
	span.SetAttributes(attribute.Int("inventory.result.parts_count", len(*parts)))

	protoParts := make([]*inventoryV1.Part, 0, len(*parts))
	for _, part := range *parts {
		protoPart := converter.ToProtoPart(&part)
		protoParts = append(protoParts, protoPart)
	}

	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
