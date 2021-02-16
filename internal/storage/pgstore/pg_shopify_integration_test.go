package pgstore

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
)

func TestAddShopifyOrderList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`insert into shopify_orders 
						(org_id,
                         shopify_account_id, 
                         order_id, 
                         created_at,
                         product_id, 
                         amount_product, 
                         bonus_id, 
                         bonus_per_piece_usd)
						values ($1,$2,$3,$4,$5,$6,$7,$8)`)).WithArgs(
		1, "12345", 1234567, "2021-02-15T22:44:11+01:00", 123456789, 2, 1, 100).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	st := toStore(db)

	err = st.AddShopifyOrderList(ctx, []user.Order{
		{
			ID:               0,
			OrganizationID:   1,
			ShopifyAccountID: "12345",
			CreatedAt:        "2021-02-15T22:44:11+01:00",
			ProductID:        123456789,
			OrderID:          1234567,
			AmountProduct:    2,
			BonusID:          1,
			BonusPerPieceUSD: 100,
		},
	})
	if err != nil {
		t.Fatalf("cannot insert new shopify order: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetOrdersByShopifyAccountID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from shopify_orders where shopify_account_id = $1`)).WithArgs(
		"12345").WillReturnRows(sqlmock.NewRows([]string{
		"id", "org_id", "shopify_account_id", "created_at",
		"product_id", "order_id", "amount_product", "bonus_id", "bonus_per_piece_usd",
	}).AddRow(0, 1, "12345", "2021-02-15T22:44:11+01:00", 123456789, 1234567, 2, 1, 100))

	st := toStore(db)

	orders, err := st.GetOrdersByShopifyAccountID(ctx, "12345")
	if err != nil {
		t.Fatalf("cannot get orders by shopify account: %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("len of items returned not matching")
	}

	if orders[0].ID != 0 || orders[0].OrganizationID != 1 || orders[0].ShopifyAccountID != "12345" || orders[0].CreatedAt != "2021-02-15T22:44:11+01:00" || orders[0].ProductID != 123456789 ||
		orders[0].OrderID != 1234567 || orders[0].AmountProduct != 2 || orders[0].BonusID != 1 || orders[0].BonusPerPieceUSD != 100 {
		t.Fatalf("item not matching")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetLastOrderByShopifyAccountID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from shopify_orders where shopify_account_id = $1 order by created_at desc offset 0 limit 1`)).WithArgs(
		"12345").WillReturnRows(sqlmock.NewRows([]string{
		"id", "org_id", "shopify_account_id", "created_at",
		"product_id", "order_id", "amount_product", "bonus_id", "bonus_per_piece_usd",
	}).AddRow(0, 1, "12345", "2021-02-15T22:44:11+01:00", 123456789, 1234567, 2, 1, 100))

	st := toStore(db)
	order, err := st.GetLastOrderByShopifyAccountID(ctx, "12345")
	if err != nil {
		t.Fatalf("cannot get orders by shopify account: %v", err)
	}

	if order.ID != 0 || order.OrganizationID != 1 || order.ShopifyAccountID != "12345" || order.CreatedAt != "2021-02-15T22:44:11+01:00" || order.ProductID != 123456789 ||
		order.OrderID != 1234567 || order.AmountProduct != 2 || order.BonusID != 1 || order.BonusPerPieceUSD != 100 {
		t.Fatalf("item not matching")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetOrdersWithPendingBonusStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from shopify_orders where bonus_id = 0 offset $1 limit $2`)).WithArgs(
		0, 11).WillReturnRows(sqlmock.NewRows([]string{
		"id", "org_id", "shopify_account_id", "created_at",
		"product_id", "order_id", "amount_product", "bonus_id", "bonus_per_piece_usd",
	}).AddRow(0, 1, "12345", "2021-02-15T22:44:11+01:00", 123456789, 1234567, 2, 1, 100))

	st := toStore(db)
	orders, err := st.GetOrdersWithPendingBonusStatus(ctx, 0, 11)
	if err != nil {
		t.Fatalf("cannot get orders by shopify account: %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("len of items returned not matching")
	}

	if orders[0].ID != 0 || orders[0].OrganizationID != 1 || orders[0].ShopifyAccountID != "12345" || orders[0].CreatedAt != "2021-02-15T22:44:11+01:00" || orders[0].ProductID != 123456789 ||
		orders[0].OrderID != 1234567 || orders[0].AmountProduct != 2 || orders[0].BonusID != 1 || orders[0].BonusPerPieceUSD != 100 {
		t.Fatalf("item not matching")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetOrdersCountWithPendingBonusStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	mock.ExpectQuery(regexp.QuoteMeta(`select count(id) from shopify_orders where bonus_id = 0`)).WillReturnRows(sqlmock.NewRows([]string{
		"count",
	}).AddRow(1))

	st := toStore(db)
	count, err := st.GetOrdersCountWithPendingBonusStatus(ctx)
	if err != nil {
		t.Fatalf("cannot get count of orders with pending bonus status: %v", err)
	}

	if count != 1 {
		t.Fatalf("not matching")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateBonusID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	mock.ExpectExec(regexp.QuoteMeta(`update shopify_orders set bonus_id = $1 where id = $2`)).WillReturnResult(sqlmock.NewResult(0, 1))

	st := toStore(db)
	err = st.UpdateBonusID(ctx, 1, 23)
	if err != nil {
		t.Fatalf("cannot update bonus id: %v", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
