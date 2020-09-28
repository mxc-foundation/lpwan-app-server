// Package congecko implements a client to coin gecko API, see details at
// https://www.coingecko.com/en/api
package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client is the object that provides access to API
type Client struct {
	c       *http.Client
	baseURL string
}

// New creates a new Client
func New() *Client {
	return &Client{
		c:       http.DefaultClient,
		baseURL: "https://api.coingecko.com/api/v3",
	}
}

// map[crypto]map[fiat]price
type simplePriceResponse map[string]map[string]float64

// GetPrice returns price of the specified crypto currency in the units of the
// specified fiat currency
func (c *Client) GetPrice(crypto, fiat string) (float64, error) {
	u := c.baseURL + "/simple/price?"
	v := url.Values{}
	v.Set("ids", crypto)
	v.Set("vs_currencies", fiat)
	u += v.Encode()
	resp, err := c.c.Get(u)
	if err != nil || resp.StatusCode != 200 {
		return 0, fmt.Errorf("couldn't get price error: %v status: %s", err, resp.Status)
	}
	dec := json.NewDecoder(resp.Body)
	var sr simplePriceResponse
	if err := dec.Decode(&sr); err != nil {
		return 0, fmt.Errorf("couldn't decode response: %v", err)
	}
	price, ok := sr[crypto][fiat]
	if !ok {
		return 0, fmt.Errorf("price is missing from response: %#v", sr)
	}
	return price, nil
}
