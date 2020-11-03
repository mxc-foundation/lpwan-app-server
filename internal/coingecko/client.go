// Package congecko implements a client to coin gecko API, see details at
// https://www.coingecko.com/en/api
package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const defaultTTL = 5 * time.Minute

type pair struct {
	crypto string
	fiat   string
}

type cachedPrice struct {
	price   float64
	expires time.Time
}

// Client is the object that provides access to API
type Client struct {
	c       *http.Client
	baseURL string
	cache   map[pair]cachedPrice
	muCache sync.RWMutex
	ttl     time.Duration
}

// New creates a new Client
func New() *Client {
	return &Client{
		c:       http.DefaultClient,
		baseURL: "https://api.coingecko.com/api/v3",
		cache:   make(map[pair]cachedPrice),
		ttl:     defaultTTL,
	}
}

// map[crypto]map[fiat]price
type simplePriceResponse map[string]map[string]float64

// GetPrice returns price of the specified crypto currency in the units of the
// specified fiat currency
func (c *Client) GetPrice(crypto, fiat string) (float64, error) {
	p := pair{crypto: crypto, fiat: fiat}
	c.muCache.RLock()
	cached, ok := c.cache[p]
	c.muCache.RUnlock()

	if !ok || cached.expires.Before(time.Now()) {
		// if we don't have cached price or if it has expired retrieve updated
		// one from coingecko and store it in the cache
		price, err := c.updatePrice(p)
		if err != nil {
			return 0, err
		}
		c.muCache.Lock()
		c.cache[p] = cachedPrice{
			price:   price,
			expires: time.Now().Add(c.ttl),
		}
		c.muCache.Unlock()
		return price, nil
	}

	return cached.price, nil
}

func (c *Client) updatePrice(p pair) (float64, error) {
	u := c.baseURL + "/simple/price?"
	v := url.Values{}
	v.Set("ids", p.crypto)
	v.Set("vs_currencies", p.fiat)
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
	price, ok := sr[p.crypto][p.fiat]
	if !ok {
		return 0, fmt.Errorf("price is missing from response: %#v", sr)
	}
	return price, nil
}
