package sqlite

import (
	"context"

	"github.com/egustafson/uberwerks/kvd-go/store"
)

const DriverID = "sqlite-store"

type sqliteStoreDriver struct{}

func init() {
	store.RegisterDriver(new(sqliteStoreDriver))
}

func (drv *sqliteStoreDriver) DriverID() string {
	return DriverID
}

func (drv *sqliteStoreDriver) New(ctx context.Context, params map[string]string) (store.Store, error) {
	return nil, nil
}
