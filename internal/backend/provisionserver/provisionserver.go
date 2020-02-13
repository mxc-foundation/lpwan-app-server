package provisionserver

import (
	"bytes"
	api "github.com/mxc-foundation/lpwan-app-server/api/ps"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"google.golang.org/grpc"
	"sync"
)

var p Pool

// Pool defines the provision-server client pool.
type Pool interface {
	Get(hostname string, caCert, tlsCert, tlsKey []byte) (api.ProvisionServiceClient, error)
}

type client struct {
	client api.ProvisionServiceClient
	clientConn *grpc.ClientConn
	caCert     []byte
	tlsCert    []byte
	tlsKey     []byte
}

// Setup configures the provision-server package.
func Setup(conf config.Config) error {
	p = &pool{
		clients: make(map[string]client),
	}

	return nil
}

type pool struct {
	sync.RWMutex
	clients map[string]client
}

func GetPool() Pool {
	return p
}

func SetPool(pool Pool) {
	p = pool
}

// Get returns a ProvisionServerClient for the given server (hostname:ip).
func (p pool) Get(hostname string, caCert, tlsCert, tlsKey []byte) (api.ProvisionServiceClient, error) {
	p.Lock()
	defer p.Unlock()

	var connect bool
	c, ok := p.clients[hostname]
	if !ok {
		connect = true
	}

	// if the connection exists in the map, but when the certificates changed
	// try to cloe the connection and re-connect
	if ok && (!bytes.Equal(c.caCert, caCert) || !bytes.Equal(c.tlsCert, tlsCert) || !bytes.Equal(c.tlsKey, tlsKey)) {
		c.clientConn.Close()
		delete(p.clients, hostname)
		connect = true
	}

	if connect {

	}
}