package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/test"
)

// DatabaseTestSuiteBase provides the setup and teardown of the database
// for every test-run.
type DatabaseTestSuiteBase struct {
	ctx context.Context
	tx  *store.Handler
}

// SetupSuite is called once before starting the test-suite.
func (b *DatabaseTestSuiteBase) SetupSuite() {
	conf := test.GetConfig()
	if err := Setup(conf); err != nil {
		panic(err)
	}
}

// SetupTest is called before every test.
func (b *DatabaseTestSuiteBase) SetupTest() {
	b.ctx = context.Background()
	handler, err := store.New(pgstore.New(DBTest().DB))
	if err != nil {
		panic(err)
	}

	tx, err := handler.TxBegin(b.ctx)
	if err != nil {
		panic(err)
	}
	b.tx = tx

	test.MustResetDB(DBTest().DB)
	RedisClient().FlushAll()
}

// TearDownTest is called after every test.
func (b *DatabaseTestSuiteBase) TearDownTest() {
	if err := b.tx.TxRollback(b.ctx); err != nil {
		panic(err)
	}
}

// Tx returns a database transaction (which is rolled back after every
// test).
func (b *DatabaseTestSuiteBase) Tx() *store.Handler {
	return b.tx
}

type StorageTestSuite struct {
	suite.Suite
	DatabaseTestSuiteBase
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
