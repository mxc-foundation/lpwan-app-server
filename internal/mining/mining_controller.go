package mining

import (
	"context"
	"encoding/json"
	"github.com/brocaar/lorawan"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_serves_appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
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
	}
}

func Setup(conf config.Config) error {
	log.Info("mining cron task begin...")
	c := cron.New()
	// everyday 3 am
	err := c.AddFunc("0 0 3 * * ?", func() {
		log.Info("Start token mining")
		go func() {
			err := tokenMining(context.Background(), conf)
			if err != nil {
				log.WithError(err).Error("tokenMining Error")
			}
		}()
	})
	if err != nil {
		log.Fatal(err)
	}

	go c.Start()

	return nil
}

func getUSDprice(conf config.Config) float64 {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/tools/price-conversion", nil)
	if err != nil {
		log.WithError(err).Error("CMC client error")
		os.Exit(1)
	}

	q := url.Values{}
	//q.Add("id", "2")
	q.Add("symbol", "USD")
	q.Add("amount", "1")
	q.Add("convert", "MXC")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", conf.ApplicationServer.MiningSetUp.CMCKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("CMC request error")
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	var cmc CMC
	err = json.Unmarshal(respBody, &cmc)
	if err != nil {
		log.Println("JSON unmarshal error: ", err)
	}

	return cmc.Data.Quote.MXC.Price
}

func tokenMining(ctx context.Context, conf config.Config) error {
	current_time := time.Now().Unix()

	mining_gws, err := storage.GetGatewayMiningList(ctx, storage.DB(), current_time)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			log.Info("No gateway online longer than 24 hours")
			return nil
		}

		log.WithError(err).Error("Cannot get mining gateway list from DB.")
		return err
	}

	if len(mining_gws) == 0 {
		return nil
	}

	// usd between 10 - 12
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 12
	randNum := float64(rand.Intn(max-min+1) + min)
	mxc_price := getUSDprice(conf)
	// amount = 1 MXC:USD * rand amount / amount of gateways
	amount := mxc_price * randNum / float64(len(mining_gws))

	var macs []string

	// update the first heartbeat = current_time
	for _, v := range mining_gws {
		err := storage.UpdateFirstHeartbeat(ctx, storage.DB(), v, current_time)
		if err != nil {
			log.WithError(err).Error("tokenMining/update first heartbeat error")
		}
		mac := lorawan.EUI64.String(v)
		macs = append(macs, mac)
	}

	miningSent := false
	// if error, resend after one minute
	for !miningSent {
		err := sendMining(ctx, macs, mxc_price, amount)
		if err != nil {
			log.WithError(err).Error("send mining request to m2m error")
			time.Sleep(60 *time.Second)
			continue
		}
		miningSent = true
	}

	return nil
}

func sendMining(ctx context.Context, macs []string, mxc_price, amount float64) error {
	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error("create m2mClient for mining error")

		return err
	}

	miningClient := api.NewMiningServiceClient(m2mClient)

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