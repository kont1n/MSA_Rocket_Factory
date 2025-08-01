package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	repoOrder := converter.ToRepoOrderPostgres(order)

	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "part_uuid", "total_price", "transaction_uuid", "payment_method", "status").
		Values(repoOrder.UserUUID, repoOrder.PartUUIDs, repoOrder.TotalPrice, repoOrder.TransactionUUID, repoOrder.PaymentMethod, repoOrder.Status).
		Suffix("RETURNING order_uuid")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, model.ErrFailedToBuildQuery
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&repoOrder.OrderUUID)
	if err != nil {
		return nil, model.ErrFailedToInsertOrder
	}

	order.OrderUUID = repoOrder.OrderUUID

	return order, nil
}
