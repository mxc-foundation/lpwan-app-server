package mining

import (
	"context"
	"fmt"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	"math/rand"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
	"github.com/mxc-foundation/lpwan-app-server/internal/coingecko"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type CMC struct {
	Status *Status `json:"status"`
	Data   *Data   `json:"data"`
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
	Notice       string `json:"notice"`
}

type Data struct {
	Id          int    `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	LastUpdated string `json:"last_updated"`
	Quote       struct {
		MXC struct {
			Price      float64 `json:"price"`
			LastUpdate string  `json:"last_update"`
		}
		USD struct {
			Price      float64 `json:"price"`
			LastUpdate string  `json:"last_update"`
		}
	}
}

var dailyMxcPrice, totalMiningValue float64
var rnd *rand.Rand

// PriceFetcher provides method to fetch current crypto currency price
type PriceFetcher interface {
	// GetPrice returns price of the specified crypto currency against
	// specified fiat currency
	GetPrice(crypto, fiat string) (float64, error)
}

type Controller struct {
	priceFetcher  PriceFetcher
	rnd           *rand.Rand
	crypto        string
	fiat          string
	minFiatValue  float64
	maxFiatValue  float64
	gwOnlineLimit int64
	lastPrice     float64
}

func Setup(conf config.Config) error {
	log.Info("mining cron task begin...")

	mconf := conf.ApplicationServer.MiningSetUp
	if mconf.MinValue <= 0 || mconf.MinValue > mconf.MaxValue {
		err := fmt.Errorf("invalid mining configuration, min_value %d and max_value %d", mconf.MinValue, mconf.MaxValue)
		log.Error(err)
		return err
	}

	ctrl := &Controller{
		priceFetcher: coingecko.New(),
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
		crypto:       "mxc",
		fiat:         "usd",
		// converting value/day to value/10 minutes
		minFiatValue:  float64(mconf.MinValue) / 144,
		maxFiatValue:  float64(mconf.MaxValue) / 144,
		gwOnlineLimit: mconf.GwOnlineLimit,
	}

	c := cron.New()
	exeTime := config.C.ApplicationServer.MiningSetUp.ExecuteTime

	err := c.AddFunc(exeTime, func() {
		log.Info("Start token mining")
		go func() {
			err := ctrl.tokenMining(context.Background(), conf)
			if err != nil {
				log.WithError(err).Error("tokenMining Error")
			}
		}()
	})
	if err != nil {
		log.WithError(err).Error("Start mining cron task failed")
	}
	go c.Start()

	return nil
}

func (ctrl *Controller) tokenMining(ctx context.Context, conf config.Config) error {
	price, err := ctrl.priceFetcher.GetPrice(ctrl.crypto, ctrl.fiat)
	if err != nil {
		log.WithError(err).Errorf("couldn't get the price of %s", ctrl.crypto)
		if ctrl.lastPrice == 0 {
			return fmt.Errorf("couldn't get the price of %s and don't have last price", ctrl.crypto)
		}
		price = ctrl.lastPrice
	}
	ctrl.lastPrice = price
	current_time := time.Now().Unix()

	// get the gateway list that should receive the mining tokens
	miningGws, err := gateway.Service.St.GetGatewayMiningList(ctx, current_time, ctrl.gwOnlineLimit)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			log.Info("No gateway online longer than 24 hours")
			return nil
		}

		log.WithError(err).Error("Cannot get mining gateway list from DB.")
		return err
	}

	if len(miningGws) == 0 {
		return nil
	}

	var macs []string

	// update the first heartbeat = 0
	for _, v := range miningGws {
		err := gateway.Service.St.UpdateFirstHeartbeatToZero(ctx, v)
		if err != nil {
			log.WithError(err).Error("tokenMining/update first heartbeat to zero error")
		}
		mac := lorawan.EUI64.String(v)
		macs = append(macs, mac)
	}

	amount := (ctrl.minFiatValue + ctrl.rnd.Float64()*
		(ctrl.maxFiatValue-ctrl.minFiatValue)) / price

	miningSent := false
	// if error, resend after one minute
	for !miningSent {
		err := sendMining(ctx, macs, 1/price, amount)
		if err != nil {
			log.WithError(err).Error("send mining request to m2m error")
			time.Sleep(60 * time.Second)
			continue
		}
		miningSent = true
	}

	return nil
}

func sendMining(ctx context.Context, macs []string, mxc_price, amount float64) error {
	miningClient, err := m2mcli.GetMiningServiceClient()
	if err != nil {
		log.WithError(err).Error("create m2mClient for mining error")
		return err
	}

	resp, err := miningClient.Mining(ctx, &api.MiningRequest{
		GatewayMac:    macs,
		MiningRevenue: amount,
		MxcPrice:      mxc_price,
	})
	if err != nil {
		log.WithError(err).Error("Mining API request error")
		return err
	}

	// if response == false, resend the request to m2m
	for !resp.Status {
		time.Sleep(60 * time.Second)
		log.Println("Resend mining request.......")
		resp, err = miningClient.Mining(ctx, &api.MiningRequest{
			GatewayMac:    macs,
			MiningRevenue: amount,
			MxcPrice:      mxc_price,
		})
		if err != nil {
			log.WithError(err).Error("Mining API request error")
			return err
		}
	}

	return nil
}
