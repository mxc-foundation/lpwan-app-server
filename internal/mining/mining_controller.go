package mining

import (
	"encoding/json"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
		MXC struct{
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
			err := tokenMining(conf)
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
	req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v1/tools/price-conversion", nil)
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

func tokenMining(conf config.Config) error {
	//t := time.Now()
	//first date of month 00:00:00
	//startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	//last date of month 23:59:59
	//endTime := startTime.AddDate(0, 1, 0).Add(time.Second * -1)
	

	//ToDo:
	// if current_heartbeat - last_heartbeat > xxx (for example three heartbeats), then first_heartbeat = current_heartbeat && last_heartbeat = current_heartbeat
	// if last_heartbeat - first_heartbeat >= 24hrs && endTime - first_heartbeat >= 24 hrs { give the tokens }

	return nil
}