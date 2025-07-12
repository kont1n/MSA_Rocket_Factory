package v1	

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (c *client) GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	return nil, nil
}
