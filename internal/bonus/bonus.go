// Package bonus implements bonus payment system
package bonus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	bapi "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
)

// Config contains configuration of the service
type Config struct {
	// URL of the server that provides the list of bonuses and credentials
	URL      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	// How often to retrieve and pay the bonuses
	CheckInterval int64 `mapstructure:"check_interval_sec"`
	// the identificator for this supernode used by remote side
	SNID string `mapstructure:"supernode_id"`
}

// Store is the DB interface
type Store interface {
	// GetUserBonusOrgID given the user's email address returns the ID of the
	// organization to which the bonuses awarded to the user should be paid. If
	// user does not exist returns 0, nil
	GetUserBonusOrgID(ctx context.Context, email string) (int64, error)
}

// Service represents an instance of the running service
type Service struct {
	httpCli  *http.Client
	cfg      Config
	jwt      string
	interval time.Duration
	store    Store
	dbCli    bapi.DistributeBonusServiceClient
	done     chan struct{}
}

// Start starts the service
func Start(ctx context.Context, cfg Config, store Store, dbCli bapi.DistributeBonusServiceClient) *Service {
	if cfg.URL == "" {
		logrus.Infof("URL for bonus client is not specified, not starting")
		return nil
	}
	srv := &Service{
		httpCli:  &http.Client{},
		cfg:      cfg,
		interval: time.Duration(cfg.CheckInterval) * time.Second,
		store:    store,
		dbCli:    dbCli,
		done:     make(chan struct{}),
	}
	go srv.run()
	return srv
}

// Stop stops the service. The service object is not usable after this call
func (srv *Service) Stop() {
	if srv != nil {
		srv.done <- struct{}{}
		close(srv.done)
	}
}

func (srv *Service) run() {
	for {
		wait := time.Until(time.Now().Truncate(srv.interval).Add(srv.interval))
		select {
		case <-time.After(wait):
			if err := srv.processAirdrops(); err != nil {
				logrus.Errorf("failed to process airdrops: %v", err)
			}
		case <-srv.done:
			return
		}
	}
}

type authReq struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type authResp struct {
	JWT string `json:"jwt"`
}

func (srv *Service) authenticate() error {
	ar := &authReq{
		Identifier: srv.cfg.User,
		Password:   srv.cfg.Password,
	}
	b, err := json.Marshal(ar)
	if err != nil {
		return err
	}
	resp, err := srv.httpCli.Post(srv.cfg.URL+"/auth/local", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("couldn't send authentication request: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("authentication has failed, status: %s", resp.Status)
	}
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't read authentication response: %v", err)
	}
	var authInfo authResp
	if err := json.Unmarshal(rb, &authInfo); err != nil {
		return fmt.Errorf("couldn't decode authentication response: %v", err)
	}
	if authInfo.JWT == "" {
		return fmt.Errorf("JWT token is empty")
	}
	srv.jwt = authInfo.JWT
	return nil
}

type airdrop struct {
	ID          int64   `json:"id"`
	Email       string  `json:"email"`
	Supernode   string  `json:"supernode"`
	Token       string  `json:"token"`
	AmountUSD   float64 `json:"amount"`
	Purpose     string  `json:"purpose"`
	Distributed bool    `json:"distributed"`
	Error       string  `json:"error"`
}

func (srv *Service) getListOfAirdrops() ([]airdrop, error) {
	v := url.Values{}
	v.Set("supernode", srv.cfg.SNID)
	v.Set("distributed", "false")
	v.Set("error", "")
	req, err := http.NewRequest("GET", srv.cfg.URL+"/airdrops?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+srv.jwt)
	resp, err := srv.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send request for airdrops: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request for airdrops returned an error: %s", resp.Status)
	}
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read list of airdrops: %v", err)
	}
	var list []airdrop
	if err := json.Unmarshal(rb, &list); err != nil {
		return nil, fmt.Errorf("couldn't decode airdrops list: %v", err)
	}
	return list, nil
}

func mapToken(token string) (string, error) {
	if token == "mxc" {
		return "ETH_MXC", nil
	} else if token == "btc" {
		return "BTC", nil
	} else if token == "dhx" {
		return "DHX", nil
	}
	return "", fmt.Errorf("unknown token: %s", token)
}

func (srv *Service) processAirdrops() error {
	if err := srv.authenticate(); err != nil {
		return err
	}
	ads, err := srv.getListOfAirdrops()
	if err != nil {
		return nil
	}
	ctx := context.Background()
	for _, ad := range ads {
		orgID, err := srv.store.GetUserBonusOrgID(ctx, ad.Email)
		if err != nil {
			return fmt.Errorf("couldn't get user's orgId: %v", err)
		}
		if orgID == 0 {
			srv.updateAirdropError(ad, "the user does not exist or doesn't have an organization")
			continue
		}
		currency, err := mapToken(ad.Token)
		if err != nil {
			srv.updateAirdropError(ad, err.Error())
			continue
		}
		req := &bapi.AddBonusRequest{
			OrgId:       orgID,
			Currency:    currency,
			AmountUsd:   fmt.Sprintf("%g", ad.AmountUSD),
			Description: ad.Purpose,
			ExternalRef: fmt.Sprintf("airdrop-%d", ad.ID),
		}
		_, err = srv.dbCli.AddBonus(ctx, req)
		if err != nil {
			srv.updateAirdropError(ad, err.Error())
			continue
		}
		srv.updateAirdropPaid(ad)
	}
	return nil
}

func (srv *Service) updateAirdrop(airdropID int64, update interface{}) error {
	b, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("couldn't marshal the update: %v", err)
	}
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/airdrops/%d", srv.cfg.URL, airdropID),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return fmt.Errorf("couldn't create update request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+srv.jwt)
	req.Header.Add("Content-Type", "application/json")
	resp, err := srv.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("couldn't send update request: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("couldn't update airdrop %d: %s", airdropID, resp.Status)
	}
	logrus.Infof("update request: %s", string(b))
	return nil
}

type errorUpdate struct {
	Error string `json:"error"`
}

func (srv *Service) updateAirdropError(ad airdrop, errStr string) {
	logrus.Infof("failed to pay airdrop %d: %s", ad.ID, errStr)
	if err := srv.updateAirdrop(ad.ID, &errorUpdate{Error: errStr}); err != nil {
		logrus.Errorf("failed to update airdrop status upstream: %v", err)
	}
}

type paidUpdate struct {
	Distributed bool `json:"distributed"`
}

func (srv *Service) updateAirdropPaid(ad airdrop) {
	logrus.Infof("successfully distributed airdrop %d", ad.ID)
	if err := srv.updateAirdrop(ad.ID, &paidUpdate{Distributed: true}); err != nil {
		logrus.Errorf("failed to update airdrop status upstream: %v", err)
	}
}
