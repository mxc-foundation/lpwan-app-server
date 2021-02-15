package bonus

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/grpc"

	bapi "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
)

func handleAuth(rw http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	var auth map[string]string
	if err := json.Unmarshal(b, &auth); err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	if auth["identifier"] != "foo" || auth["password"] != "boo" {
		err := fmt.Errorf("no auth info in request: %#v", auth)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte(`{"jwt":"bar"}`))
}

func authenticated(r *http.Request) error {
	auth := r.Header.Get("Authorization")
	if auth != `Bearer bar` {
		return fmt.Errorf("not authenticated")
	}
	return nil
}

var airdrops = []*airdrop{
	{
		ID:    1,
		Email: "foo@example.com",
	},
	{
		ID:    2,
		Email: "boo@example.com",
	},
}

func handleAirdrops(rw http.ResponseWriter, r *http.Request) {
	if err := authenticated(r); err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}
	if r.Method == "GET" {
		handleList(rw, r)
	} else if r.Method == "PUT" {
		handleUpdate(rw, r)
	} else {
		http.Error(rw, "unknown method "+r.Method, http.StatusMethodNotAllowed)
	}
}

func handleList(rw http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	if vals.Get("supernode") != "unit" || vals.Get("distributed") != "false" || vals["error"] == nil || vals.Get("error") != "" {
		http.Error(rw, "some parameters are missing", http.StatusBadRequest)
		return
	}
	var res []*airdrop
	for _, ad := range airdrops {
		ad.Supernode = "unit"
		ad.Token = "mxc"
		ad.AmountUSD = 10
		ad.Purpose = "test"
		ad.Distributed = false
		ad.Error = ""
		res = append(res, ad)
	}
	b, err := json.Marshal(res)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(b)
}

func handleUpdate(rw http.ResponseWriter, r *http.Request) {
	var id int64
	if r.URL.Path == "/airdrops/1" {
		id = 1
	} else if r.URL.Path == "/airdrops/2" {
		id = 2
	} else {
		err := fmt.Errorf(r.URL.Path + " not found")
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	var req map[string]interface{}
	if err := json.Unmarshal(b, &req); err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	if d, ok := req["distributed"].(bool); ok {
		airdrops[id-1].Distributed = d
	}
	if e, ok := req["error"].(string); ok {
		airdrops[id-1].Error = e
	}
	b, err = json.Marshal(airdrops[id-1])
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(b)
}

type testStore struct{}

func (ts *testStore) GetUserBonusOrgID(ctx context.Context, email string) (int64, error) {
	if email == "foo@example.com" {
		return 4, nil
	}
	return 0, nil
}

type testM2M struct{}

func (tm *testM2M) AddBonus(ctx context.Context, in *bapi.AddBonusRequest, opts ...grpc.CallOption) (*bapi.AddBonusResponse, error) {
	return &bapi.AddBonusResponse{}, nil
}

func TestProcessAirdrops(t *testing.T) {
	m2m := &testM2M{}
	ts := &testStore{}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/local", handleAuth)
	mux.HandleFunc("/airdrops/", handleAirdrops)
	hs := httptest.NewServer(mux)
	defer hs.Close()

	srv := &Service{
		httpCli: &http.Client{},
		cfg: Config{
			URL:      hs.URL,
			User:     "foo",
			Password: "boo",
			SNID:     "unit",
		},
		interval: 3600 * time.Second,
		store:    ts,
		dbCli:    m2m,
		done:     make(chan struct{}),
	}
	if err := srv.processAirdrops(); err != nil {
		t.Errorf("processAirdrops returned an error: %v", err)
	}
	if !airdrops[0].Distributed || airdrops[0].Error != "" {
		t.Errorf("airdrop 1 is not as expected %#v", airdrops[0])
	}
	if airdrops[1].Distributed || airdrops[1].Error == "" {
		t.Errorf("airdrop 2 is not as expected %#v", airdrops[1])
	}
}
