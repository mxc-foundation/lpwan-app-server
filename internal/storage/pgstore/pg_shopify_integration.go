package pgstore

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
)

// AddShopifyOrder inserts new order record into db
func (ps *PgStore) AddShopifyOrder(ctx context.Context, order user.Order) error {
	_, err := ps.db.ExecContext(ctx,
		`insert into shopify_orders (user_id, order_id, created_at, product_id, order_amount, bonus_id)
				values ($1,$2,$3,$4,$5,$6)`,
		order.UserID,
		order.OrderID,
		order.CreatedAt,
		order.ProductID,
		order.Amount,
		order.BonusID)
	return err
}

// GetOrdersByUserID returns order list with given user id
func (ps *PgStore) GetOrdersByUserID(ctx context.Context, userID int64) ([]user.Order, error) {
	var orders []user.Order

	err := sqlx.SelectContext(ctx, ps.db, &orders, `select * from shopify_orders where user_id = $1`, userID)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return orders, nil
}

// GetLastOrderByUserID returns last order information for the given user id
func (ps *PgStore) GetLastOrderByUserID(ctx context.Context, userID int64) (user.Order, error) {
	var order user.Order

	err := sqlx.GetContext(ctx, ps.db, &order, `
			select * from shopify_orders where user_id = $1 order by created_at desc offset 0 limit 1`, userID)
	if err != nil {
		return order, handlePSQLError(Select, err, "select error")
	}
	return order, nil
}
