package coingecko

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPrice(t *testing.T) {
	testPrices := map[string]map[string]float64{
		"mxc": {"usd": 0.01199642},
		"eth": {"eur": 227.39},
	}
	var requestsCnt int
	m := http.NewServeMux()
	m.HandleFunc("/simple/price", func(w http.ResponseWriter, r *http.Request) {
		requestsCnt++
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		id := r.Form.Get("ids")
		vsc := r.Form.Get("vs_currencies")
		price, ok := testPrices[id][vsc]
		if ok {
			fmt.Fprintf(w, `{"%s":{"%s":%.9f}}`, id, vsc, price)
		} else {
			fmt.Fprintf(w, `{}`)
		}
	})
	s := httptest.NewServer(m)
	defer s.Close()

	c := New()
	c.baseURL = s.URL

	mup, err := c.GetPrice("mxc", "usd")
	if err != nil {
		t.Errorf("got an error for mxc/usd: %v", err)
	}
	if math.Abs(mup-testPrices["mxc"]["usd"]) > 1e-9 {
		t.Errorf("wrong mxc/usd price: %f", mup)
	}
	// this request should use cached price
	_, _ = c.GetPrice("mxc", "usd")
	if requestsCnt > 1 {
		t.Errorf("expected second mxc/usd request to use cached data")
	}

	_, err = c.GetPrice("mxc", "eur")
	if err == nil || !strings.Contains(err.Error(), "price is missing") {
		t.Errorf("unexpected error for mxc/eur: %v", err)
	}
}
