package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/model"
)

func (r *repository) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	builderSelect := sq.Select(
		"order_uuid", "user_uuid", "part_uuid", "total_price",
		"transaction_uuid", "payment_method", "status", "created_at", "updated_at").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_uuid": id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, model.ErrFailedToBuildQuery
	}

	var repoOrder repoModel.OrderPostgres
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&repoOrder.OrderUUID,
		&repoOrder.UserUUID,
		&repoOrder.PartUUIDs,
		&repoOrder.TotalPrice,
		&repoOrder.TransactionUUID,
		&repoOrder.PaymentMethod,
		&repoOrder.Status,
		&repoOrder.CreatedAt,
		&repoOrder.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrOrderNotFound
		}
		return nil, model.ErrFailedToGetOrder
	}

	order, err := converter.ToModelOrderFromPostgres(&repoOrder)
	if err != nil {
		return nil, err
	}

	return order, nil
}
