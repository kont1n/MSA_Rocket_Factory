package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	repoOrder := converter.ToRepoOrderPostgres(order)

	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("transaction_uuid", repoOrder.TransactionUUID).
		Set("payment_method", repoOrder.PaymentMethod).
		Set("status", repoOrder.Status).
		Set("updated_at", "NOW()").
		Where(sq.Eq{"order_uuid": order.OrderUUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, model.ErrFailedToBuildQuery
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, model.ErrFailedToUpdateOrder
	}

	return order, nil
}
