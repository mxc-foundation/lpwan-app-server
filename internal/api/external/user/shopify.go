package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
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
	store    Store
	bonusCli pb.DistributeBonusServiceClient
	auth     auth.Authenticator
}

// NewShopifyServiceServer creates a new shopify integration service server
func NewShopifyServiceServer(cli pb.DistributeBonusServiceClient, auth auth.Authenticator, store Store) *ShopifyServiceServer {
	return &ShopifyServiceServer{
		bonusCli: cli,
		auth:     auth,
		store:    store,
	}
}

// ShopifyStore defines db APIs for dhx service
type ShopifyStore interface {
	// AddShopifyOrder inserts new order record into db
	AddShopifyOrder(ctx context.Context, order Order) error
	// GetOrdersByUserID returns order list with given user id
	GetOrdersByUserID(ctx context.Context, userID int64) ([]Order, error)
	// GetLastOrderByUserID returns last order information for the given user id
	GetLastOrderByUserID(ctx context.Context, userID int64) (Order, error)
}

// Order represent db data in table shopify_orders
type Order struct {
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ProductID string    `db:"product_id"`
	OrderID   string    `db:"order_id"`
	Amount    int64     `db:"order_amount"`
	BonusID   int64     `db:"bonus_id"`
}

// GetOrdersByUserID returns a list of shopify orders filtered by given user id
func (s *ShopifyServiceServer) GetOrdersByUserID(ctx context.Context, req *api.GetOrdersByUserIDRequest) (*api.GetOrdersByUserIDResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	extUser, err := s.store.GetExternalUserByUserIDAndService(ctx, auth.SHOPIFY, cred.UserID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to get shopify account from given user: %v", err)
	}

	orderList, err := s.store.GetOrdersByUserID(ctx, cred.UserID)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			return &api.GetOrdersByUserIDResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to get order list: %v", err)
	}

	var orders []*api.ShopifySuccessfulOrder
	for _, v := range orderList {
		orderItem := api.ShopifySuccessfulOrder{}

		orderItem.Amount = v.Amount
		orderItem.ProductId = v.ProductID
		orderItem.ShopifyAccount = extUser.ExternalUsername
		orderItem.OrderId = v.OrderID

		orderItem.CreatedAt, err = ptypes.TimestampProto(v.CreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}

		if v.BonusID == 0 {
			orderItem.BonusStatus = "processing"
		} else {
			orderItem.BonusStatus = "distributed"
		}

		orders = append(orders, &orderItem)
	}

	return &api.GetOrdersByUserIDResponse{SuccessfulOrders: orders}, nil
}

type order struct {
}

type orders struct {
	list []order
}

// Client defines client with functions that interact with shopify server
type Client struct {
	UserID   int64
	AdminAPI ShopifyAdminAPI
	Store    Store

	done chan struct{}
}

// monitoring starts go rountine for a specific user who has bound shopify account
func monitoring(ctx context.Context, userID int64, conf ShopifyAdminAPI, store Store) {
	cli := &Client{
		UserID:   userID,
		AdminAPI: conf,
		Store:    store,
	}

	go cli.run(ctx)
}

// monitoring checks whether there is new shopify order with give user id meeting requirement
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
		return time.Now().Add(1 * time.Minute), err
	}

	// check order every 12 hours
	next := time.Now().Add(12 * time.Hour)
	if time.Now().After(next) {
		// get user's last order
		lastOrder, err := cli.Store.GetLastOrderByUserID(ctx, cli.UserID)
		if err != nil {
			if err == errHandler.ErrDoesNotExist {
				// when there is no last order record, try to get full order list from shopify for this user
				url := fmt.Sprintf("https://%s:%s@%s/admin/api/%s/customers/%s/orders.json",
					cli.AdminAPI.APIKey, cli.AdminAPI.Secret, cli.AdminAPI.Hostname,
					cli.AdminAPI.APIVersion, extUser.ExternalUserID)

				var orderList orders
				if err := auth.GetHTTPResponse(url, &orderList, false); err != nil {
					// try again in 1 min
					log.Errorf("failed to get external user from user id: %v", err)
					return time.Now().Add(1 * time.Minute), err
				}

				if len(orderList.list) == 0 {
					return next, nil
				}
				// process orders, and distribute bonus
				// TODO
			}
			// something is wrong, try again in 1 min
			log.Errorf("failed to get last order with given user id: %v", err)
			return time.Now().Add(1 * time.Minute), err
		}
		// get new orders generated after last order was processed
		// TODO
		timeMin, err := ptypes.TimestampProto(lastOrder.CreatedAt)
		if err != nil {
			// something is wrong, try again in 1 min
			log.Errorf("failed to get last order's processed time: %v", err)
			return time.Now().Add(1 * time.Minute), err
		}
		url := fmt.Sprintf("https://%s:%s@%s/admin/api/%s/orders.json?status=open&created_at_min=%s",
			cli.AdminAPI.APIKey, cli.AdminAPI.Secret, cli.AdminAPI.Hostname,
			cli.AdminAPI.APIVersion, timeMin.String())

		var orderList orders
		if err := auth.GetHTTPResponse(url, &orderList, false); err != nil {
			// try again in 1 min
			log.Errorf("failed to get external user from user id: %v", err)
			return time.Now().Add(1 * time.Minute), err
		}

		if len(orderList.list) == 0 {
			return next, nil
		}
		// process orders, and distribute bonus
		// TODO

	}
	return next, nil
}

// Close terminates this goroutine
func (cli *Client) Close() error {
	cli.done <- struct{}{}
	close(cli.done)
	return nil
}
