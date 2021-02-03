package shopify

import (
	"context"
	"strconv"
	"time"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"

	log "github.com/sirupsen/logrus"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
)

type controller struct {
	config   user.Shopify
	store    user.Store
	bonusCli api.DistributeBonusServiceClient

	done chan struct{}
}

// Start initiate goroutine to process shopify order records in the background
func Start(config user.Shopify, store user.Store, cli api.DistributeBonusServiceClient) error {
	if config.Bonus.Enable == false {
		return nil
	}

	ctrl := &controller{
		config:   config,
		store:    store,
		bonusCli: cli,
	}

	go ctrl.run(context.Background())

	return nil
}

func (c *controller) processOrders(ctx context.Context) error {
	count, err := c.store.GetOrdersCountWithPendingBonusStatus(ctx)
	if err != nil {
		return err
	}

	offset := int64(0)
	limit := int64(50)

	for i := int64(0); i < count/limit; i++ {
		orderList, err := c.store.GetOrdersWithPendingBonusStatus(ctx, offset, limit)
		if err != nil {
			if err == errHandler.ErrDoesNotExist {
				return nil
			}
			return err
		}

		for _, v := range orderList {
			if v.CreatedAt.AddDate(0, 0, 30).After(time.Now()) {
				// order is not more than 30 days old, do not distribute bonus
				continue
			}

			// TODO: distribute BTC on activate of gateways

			res, err := c.bonusCli.AddBonus(ctx, &api.AddBonusRequest{
				OrgId:       v.OrganizationID,
				Currency:    "BTC",
				AmountUsd:   strconv.FormatInt(v.BonusPerPieceUSD, 64),
				Description: "Purchase M2 Pro – LPWAN Crypto-Miner from m2prominer.com ",
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

		offset = limit * (i + 1)
	}

	return nil
}

func (c *controller) nextRun(ctx context.Context) (time.Time, error) {
	// check order every 24 hours
	next := time.Now().Add(24 * time.Hour)
	if time.Now().After(next) {
		if err := c.processOrders(ctx); err != nil {
			return time.Now().Add(10 * time.Minute), nil
		}
		return next, nil
	}
	return next, nil
}

func (c *controller) run(ctx context.Context) {
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

// Close terminates this goroutine
func (c *controller) Close() error {
	c.done <- struct{}{}
	close(c.done)
	return nil
}
