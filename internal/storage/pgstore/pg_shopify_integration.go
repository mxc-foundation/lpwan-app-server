package pgstore

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
)

// AddShopifyOrderList inserts new order record into db
func (ps *PgStore) AddShopifyOrderList(ctx context.Context, orderList []user.Order) error {
	if err := ps.Tx(ctx, func(ctx context.Context, store *PgStore) error {
		for _, v := range orderList {
			_, err := ps.db.ExecContext(ctx,
				`insert into shopify_orders 
						(org_id,
                         shopify_account_id, 
                         order_id, 
                         created_at,
                         product_id, 
                         amount_product, 
                         bonus_id, 
                         bonus_per_piece_usd)
						values ($1,$2,$3,$4,$5,$6,$7,$8)`,
				v.OrganizationID,
				v.ShopifyAccountID,
				v.OrderID,
				v.CreatedAt,
				v.ProductID,
				v.AmountProduct,
				v.BonusID,
				v.BonusPerPieceUSD)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetOrdersByShopifyAccountID returns order list with given user id
func (ps *PgStore) GetOrdersByShopifyAccountID(ctx context.Context, saID string) ([]user.Order, error) {
	var orders []user.Order

	err := sqlx.SelectContext(ctx, ps.db, &orders, `select * from shopify_orders where shopify_account_id = $1`, saID)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return orders, nil
}

// GetLastOrderByShopifyAccountID returns last order information for the given user id
func (ps *PgStore) GetLastOrderByShopifyAccountID(ctx context.Context, saID string) (user.Order, error) {
	var order user.Order

	err := sqlx.GetContext(ctx, ps.db, &order, `
			select * from shopify_orders where shopify_account_id = $1 order by created_at desc offset 0 limit 1`, saID)
	if err != nil {
		return order, handlePSQLError(Select, err, "select error")
	}
	return order, nil
}

// GetOrdersWithPendingBonusStatus returns order list with bonus_id == 0
func (ps *PgStore) GetOrdersWithPendingBonusStatus(ctx context.Context, offset, limit int64) ([]user.Order, error) {
	var orderList []user.Order

	err := sqlx.SelectContext(ctx, ps.db, &orderList, `
		select * from shopify_orders where bonus_id = 0 offset $1 limit $2
	`, offset, limit)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return orderList, err
}

// GetOrdersCountWithPendingBonusStatus returns number of orders with bonus_id == 0
func (ps *PgStore) GetOrdersCountWithPendingBonusStatus(ctx context.Context) (int64, error) {
	var count int64

	err := sqlx.GetContext(ctx, ps.db, &count, `
		select count(id) from shopify_orders where bonus_id = 0 
	`)
	return count, err
}

// UpdateBonusID updates bonus_id of record in shopify_orders
func (ps *PgStore) UpdateBonusID(ctx context.Context, keyID int64, bonusID int64) error {
	_, err := ps.db.ExecContext(ctx, `
		update shopify_orders set bonus_id = $1 where id = $2
	`, bonusID, keyID)
	return err
}
