package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// ShopifyAdminAPI defines shopify admin api configuration
type ShopifyAdminAPI struct {
	Hostname   string `mapstructure:"hostname"`
	APIKey     string `mapstructure:"api_key"`
	Secret     string `mapstructure:"secret"`
	APIVersion string `mapstructure:"api_version"`
	StoreName  string `mapstructure:"store_name"`
}

// BonusSettings defines settings of shopify promitions
type BonusSettings struct {
	Enable    bool  `mapstructure:"enable"`
	ValueUSD  int64 `mapstructure:"value_usd"`
	ProductID int64 `mapstructure:"product_id"`
}

type Shopify struct {
	AdminAPI ShopifyAdminAPI `mapstructure:"shopify_admin_api"`
	Bonus    BonusSettings   `mapstructure:"bonus"`
}

// ShopifyCustomer includes a part of attributes of customer
type ShopifyCustomer struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	OrdersCount int    `json:"orders_count"`
	State       string `json:"state"`
	LastOrderID int64  `json:"last_order_id"`
}

// ShopifyCustomerList maps response of api
// https://apikey:secret@{hostname}/admin/api/2021-01/customers/search.json\?query\=email:{email}
type ShopifyCustomerList struct {
	Customers []ShopifyCustomer `json:"customers"`
}

// ShopifyServiceServer defines the Shopify integration service Server API structure
type ShopifyServiceServer struct {
	store Store
	auth  auth.Authenticator
}

// NewShopifyServiceServer creates a new shopify integration service server
func NewShopifyServiceServer(auth auth.Authenticator, store Store) *ShopifyServiceServer {
	return &ShopifyServiceServer{
		auth:  auth,
		store: store,
	}
}

// ShopifyStore defines db APIs for dhx service
type ShopifyStore interface {
	// AddShopifyOrderList inserts new order record into db
	AddShopifyOrderList(ctx context.Context, orderList []Order) error
	// GetOrdersByShopifyAccountID returns order list with given user id
	GetOrdersByShopifyAccountID(ctx context.Context, saID string) ([]Order, error)
	// GetLastOrderByShopifyAccountID returns last order information for the given user id
	GetLastOrderByShopifyAccountID(ctx context.Context, saID string) (Order, error)
	// GetOrdersWithPendingBonusStatus returns order list with bonus_id == 0
	GetOrdersWithPendingBonusStatus(ctx context.Context, offset, limit int64) ([]Order, error)
	// GetOrdersCountWithPendingBonusStatus returns number of orders with bonus_id == 0
	GetOrdersCountWithPendingBonusStatus(ctx context.Context) (int64, error)
	// UpdateBonusID updates bonus_id of record in shopify_orders
	UpdateBonusID(ctx context.Context, keyID int64, bonusID int64) error
}

// Order represent db data in table shopify_orders
type Order struct {
	ID               int64     `db:"id"`
	OrganizationID   int64     `db:"org_id"`
	ShopifyAccountID string    `db:"shopify_account_id"`
	CreatedAt        time.Time `db:"created_at"`
	ProductID        int64     `db:"product_id"`
	OrderID          int64     `db:"order_id"`
	AmountProduct    int64     `db:"amount_product"`
	BonusID          int64     `db:"bonus_id"`
	BonusPerPieceUSD int64     `db:"bonus_per_piece_usd"`
}

// GetOrdersByUser returns a list of shopify orders filtered by given email, this API is only open for global admin user
func (s *ShopifyServiceServer) GetOrdersByUser(ctx context.Context, req *api.GetOrdersByUserRequest) (*api.GetOrdersByUserResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	u, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user with given email not found: %v", err)
	}

	extUser, err := s.store.GetExternalUserByUserIDAndService(ctx, auth.SHOPIFY, u.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to get shopify account from given user: %v", err)
	}

	orderList, err := s.store.GetOrdersByShopifyAccountID(ctx, extUser.ExternalUserID)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			return &api.GetOrdersByUserResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to get order list: %v", err)
	}

	var orders []*api.Order
	for _, v := range orderList {
		orderItem := api.Order{}

		orderItem.AmountProduct = v.AmountProduct
		orderItem.ProductId = strconv.FormatInt(v.ProductID, 10)
		orderItem.ShopifyAccount = extUser.ExternalUsername
		orderItem.OrderId = strconv.FormatInt(v.OrderID, 10)

		orderItem.CreatedAt, err = ptypes.TimestampProto(v.CreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}

		if v.BonusID == 0 {
			orderItem.BonusStatus = "done"
		} else {
			orderItem.BonusStatus = "pending"
		}

		orderItem.BonusPerPieceUsd = strconv.FormatInt(v.BonusPerPieceUSD, 10)

		orders = append(orders, &orderItem)
	}

	return &api.GetOrdersByUserResponse{Orders: orders}, nil
}

// ShopifyOrder defines a part of attributes from order item
type ShopifyOrder struct {
	ID int64 `json:"id"`

	CreatedAt time.Time `json:"created_at"`

	FinancialStatus string `json:"financial_status"`

	LineItems []struct {
		ProductID int64  `json:"id"`
		Name      string `json:"name"`
	} `json:"line_items"`

	Refunds []struct {
		RefundLineItems []struct {
			LineItem struct {
				Quantity  int64 `json:"quantity"`
				ProductID int64 `json:"product_id"`
			} `json:"line_item"`
		} `json:"refund_line_items"`
	} `json:"refunds"`
}

// GetOrdersResponse defines structure of response for
//  /admin/api/2021-01/customers/{customer_id}}/orders.json?status=any
type GetOrdersResponse struct {
	Orders []ShopifyOrder `json:"orders"`
}

// Client defines client with functions that interact with shopify server
type Client struct {
	OrganizationID int64
	UserID         int64
	config         Shopify
	Store          Store

	done chan struct{}
}

// monitoring starts go rountine for a specific user who has bound shopify account
func monitoring(ctx context.Context, organizationID, userID int64, conf Shopify, store Store) {
	cli := &Client{
		OrganizationID: organizationID,
		UserID:         userID,
		config:         conf,
		Store:          store,
	}

	go cli.run(ctx)
}

func (cli *Client) parseOrders(ctx context.Context, orders []ShopifyOrder, shopifyAccountID string) error {
	var orderItem Order
	var orderList []Order

	for _, v := range orders {
		if v.FinancialStatus != "paid" {
			// skip any other orders
			continue
		}

		orderItem.OrganizationID = cli.OrganizationID
		orderItem.OrderID = v.ID
		orderItem.ShopifyAccountID = shopifyAccountID
		orderItem.ProductID = cli.config.Bonus.ProductID
		orderItem.CreatedAt = v.CreatedAt

		// parse line items
		for _, item := range v.LineItems {
			if item.ProductID != orderItem.ProductID {
				// only grant bonus with specific product
				continue
			}
			orderItem.AmountProduct += 1
		}

		// parse refunds
		for _, refund := range v.Refunds {
			for _, refunLineItem := range refund.RefundLineItems {
				if refunLineItem.LineItem.ProductID == orderItem.ProductID {
					orderItem.AmountProduct -= refunLineItem.LineItem.Quantity
				}
			}
		}

		orderItem.BonusPerPieceUSD = cli.config.Bonus.ValueUSD

		orderList = append(orderList, orderItem)
	}

	if err := cli.Store.AddShopifyOrderList(ctx, orderList); err != nil {
		return err
	}

	return nil
}

func (cli *Client) getNewOrdersFromShopify(ctx context.Context, shopifyAccount ExternalUser) error {
	// get user's last order
	lastOrder, err := cli.Store.GetLastOrderByShopifyAccountID(ctx, shopifyAccount.ExternalUserID)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			// when there is no last order record, try to get full order list from shopify for this user
			url := fmt.Sprintf("https://%s:%s@%s/admin/api/%s/customers/%s/orders.json?status=any",
				cli.config.AdminAPI.APIKey, cli.config.AdminAPI.Secret, cli.config.AdminAPI.Hostname,
				cli.config.AdminAPI.APIVersion, shopifyAccount.ExternalUserID)

			var orderList GetOrdersResponse
			if err := auth.GetHTTPResponse(url, &orderList, false); err != nil {
				// try again in 1 min
				log.Errorf("failed to get external user from user id: %v", err)
				return err
			}
			if len(orderList.Orders) == 0 {
				return nil
			}
			// if there is new order, save it
			if err := cli.parseOrders(ctx, orderList.Orders, shopifyAccount.ExternalUserID); err != nil {
				log.Errorf("failed to save orders: %v", err)
			}
			return nil
		}
		// something is wrong, try again in 1 min
		log.Errorf("failed to get last order with given user id: %v", err)
		return err
	}
	// get new orders generated after last order was processed
	timeMin, err := ptypes.TimestampProto(lastOrder.CreatedAt)
	if err != nil {
		// something is wrong, try again in 1 min
		log.Errorf("failed to get last order's processed time: %v", err)
		return err
	}

	url := fmt.Sprintf("https://%s:%s@%s/admin/api/%s/customers/%s/orders.json?status=any&created_at_min=\"%s\"",
		cli.config.AdminAPI.APIKey, cli.config.AdminAPI.Secret, cli.config.AdminAPI.Hostname,
		cli.config.AdminAPI.APIVersion, shopifyAccount.ExternalUserID, timeMin.String())
	var orderList GetOrdersResponse
	if err := auth.GetHTTPResponse(url, &orderList, false); err != nil {
		// try again in 1 min
		log.Errorf("failed to get external user from user id: %v", err)
		return err
	}
	if len(orderList.Orders) == 0 {
		return nil
	}
	// if there is new order, save it
	if err := cli.parseOrders(ctx, orderList.Orders, shopifyAccount.ExternalUserID); err != nil {
		log.Errorf("failed to save orders: %v", err)
	}
	return nil
}

func (cli *Client) nextRun(ctx context.Context) (time.Time, error) {
	// first check whether the user is still binding shopify to supernode
	extUser, err := cli.Store.GetExternalUserByUserIDAndService(ctx, auth.SHOPIFY, cli.UserID)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			// user no longer binding to shopify account, stop this go routine
			return time.Now(), nil
		}
		// try again in 1 min
		log.Errorf("failed to get external user from user id: %v", err)
		return time.Now().Add(10 * time.Minute), err
	}

	// check order every 24 hours
	next := time.Now().Add(24 * time.Hour)
	if time.Now().After(next) {
		if err := cli.getNewOrdersFromShopify(ctx, extUser); err != nil {
			return time.Now().Add(10 * time.Minute), err
		}
		return next, nil
	}
	return next, nil
}

func (cli *Client) run(ctx context.Context) {
	for {
		next, err := cli.nextRun(ctx)
		if err != nil {
			log.Errorf("check shopify order failure: %v", err)
		}
		delay := time.Until(next)
		if delay < 0 {
			// exit this goroutine
			return
		}
		select {
		case <-cli.done:
			return
		case <-time.After(delay):
		}
	}
}

// Close terminates this goroutine
func (cli *Client) Close() error {
	cli.done <- struct{}{}
	close(cli.done)
	return nil
}
