package shopify

import (
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// Service represents an running instance of the serivce
type Service struct {
	config   user.Shopify
	store    user.Store
	bonusCli api.DistributeBonusServiceClient

	done chan struct{}
}

// Start initiate goroutine to process shopify order records in the background
func Start(ctx context.Context, config user.Shopify, store user.Store, cli api.DistributeBonusServiceClient) (*Service, error) {
	if !config.Bonus.Enable {
		return nil, nil
	}

	service := &Service{
		config:   config,
		store:    store,
		bonusCli: cli,
		done:     make(chan struct{}),
	}

	go service.startCheckingOrdersForShopifyAccounts(ctx)
	go service.processOrders(ctx)

	return service, nil
}

func (c *Service) startCheckingOrdersForShopifyAccounts(ctx context.Context) {
	userAmount, err := c.store.GetUserCount(ctx)
	if err != nil {
		log.Errorf("[shopify service] failed to get user amount: %v", err)
	}

	limit := int64(50)
	for offset := int64(0); offset <= userAmount/limit; offset += limit {
		users, err := c.store.GetUsers(ctx, limit, offset)
		if err != nil {
			log.Errorf("[shopify service] failed to get users: %v", err)
		}

		for _, v := range users {
			extUser, err := c.store.GetExternalUserByUserIDAndService(ctx, auth.SHOPIFY, v.ID)
			if err != nil {
				if err == errHandler.ErrDoesNotExist {
					continue
				}
				log.Errorf("[shopify service] failed to get external user: %v", err)
			}

			if extUser.ExternalUserID == "" {
				continue
			}

			orgList, err := c.store.GetUserOrganizations(ctx, v.ID)
			if err != nil {
				log.Errorf("[shopify service] failed to get organization list: %v", err)
			}

			for _, org := range orgList {
				if !org.IsOrgAdmin {
					continue
				}
				user.CheckNewOrders(ctx, org.OrganizationID, v.ID, c.config, c.store)
			}
		}
	}
}

func (c *Service) distributeBonus(ctx context.Context) error {
	count, err := c.store.GetOrdersCountWithPendingBonusStatus(ctx)
	if err != nil {
		return err
	}

	limit := int64(50)
	for offset := int64(0); offset <= count/limit; offset += limit {
		orderList, err := c.store.GetOrdersWithPendingBonusStatus(ctx, offset, limit)
		if err != nil {
			if err == errHandler.ErrDoesNotExist {
				return nil
			}
			return err
		}

		for _, v := range orderList {
			ct, err := time.Parse("2006-01-02T15:04:05Z07:00", v.CreatedAt)
			if err != nil {
				log.Errorf("failed to convert createdAt: %v", err)
				continue
			}

			if ct.AddDate(0, 0, 30).After(time.Now()) {
				// order is not more than 30 days old, do not distribute bonus
				continue
			}

			// TODO: distribute BTC on activate of gateways
			//  Distribute bonus upon gateway activation is not implementable without matching gateway's unique information,
			// 	considering following cases :
			//   	case 1: A bought 5 from shopify. If A gifted B 2 of gateways before registration,
			//   			then A registered 3 with his organization, B registered 2 with his organizaiton,
			//  			for A we detect he has ordered 5 from shopify,
			// 				if he has coincidentally registered 2 before (not ordered from shopify,
			//				but register them after the time he placed order) ,
			//				then we detect he activated 5, number matches, we distribute 5 * bonus_per_piece
			//      case 2: A bought 5 from shopify. If A gifted B 2 of gateways before registration,
			//     			then A registered 3 with his organization, B registered 2 with his organizaiton,
			//    			for A we detect he has ordered 5, however he only registered 3 in his organization,
			//   			if we by default think all the gateways he registered are newly purchased from shopify,
			//  			we then distribute 3 * bonus_per_piece ( 2 times bonus_per_piece short )
			//      case 3: A bought 5 from shopify. A registered 3 with his organization, 2 left unopened.
			//     			We detect he ordered 5 gateways, and activated 3,
			//      		then we only distribute 3 * bonus_per_piece, however he decided to return 2,
			//     			based on refund policy, user will get 200 USD (2 * bonus_per_piece) less which is not correct

			res, err := c.bonusCli.AddBonus(ctx, &api.AddBonusRequest{
				OrgId:       v.OrganizationID,
				Currency:    "BTC",
				AmountUsd:   strconv.FormatInt(v.BonusPerPieceUSD, 10),
				Description: "Purchase M2 Pro – LPWAN Crypto-Miner from m2prominer.com ",
				ExternalRef: fmt.Sprintf("m2prominer.com-%d", v.OrderID),
			})
			if err != nil {
				log.Errorf("failed to send add bonus request to mxp: %v", err)
				continue
			}

			if err = c.store.UpdateBonusID(ctx, v.ID, res.BonusId); err != nil {
				log.Errorf("failed to update ")
				continue
			}
		}
	}

	return nil
}

func (c *Service) nextRun(ctx context.Context) (time.Time, error) {
	if err := c.distributeBonus(ctx); err != nil {
		return time.Now().Add(10 * time.Minute), nil
	}
	// check order every 24 hours
	return time.Now().Add(24 * time.Hour), nil
}

func (c *Service) processOrders(ctx context.Context) {
	for {
		next, err := c.nextRun(ctx)
		if err != nil {
			log.Errorf("process distribute bonus request error: %v", err)
		}
		delay := time.Until(next)
		select {
		case <-c.done:
			return
		case <-time.After(delay):
		}
	}
}

// Stop stops the service. The service object is not usable after this call
func (c *Service) Stop() {
	if c != nil {
		c.done <- struct{}{}
		close(c.done)
	}
}
