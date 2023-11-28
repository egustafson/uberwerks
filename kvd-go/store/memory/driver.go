package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/egustafson/uberwerks/kvd-go/store"
)

const DriverID = "mem-store"

type memStoreDriver struct{}

func init() {
	store.RegisterDriver(new(memStoreDriver))
}

func (drv *memStoreDriver) DriverID() string {
	return DriverID
}

func (drv *memStoreDriver) New(ctx context.Context, params map[string]string) (store.Store, error) {
	id, ok := params[store.StoreIdParam]
	if !ok {
		id = uuid.New().String()
	}
	store := newMemStore(ctx, id)
	return store, nil
}
