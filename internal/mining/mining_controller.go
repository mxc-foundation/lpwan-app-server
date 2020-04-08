package mining

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/robfig/cron"
)

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