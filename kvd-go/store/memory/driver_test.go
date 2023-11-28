package memory

import (
	"context"
	"testing"

	"github.com/egustafson/uberwerks/kvd-go/store"
	"github.com/stretchr/testify/assert"
)

func TestDriverID(t *testing.T) {

	driver := new(memStoreDriver)
	assert.Equal(t, DriverID, driver.DriverID())
}

func TestNew_withStoreID(t *testing.T) {

	storeID := "test-store-id"
	params := map[string]string{
		store.DriverIdParam: DriverID,
		store.StoreIdParam:  storeID,
	}
	ctx := context.Background()
	store, err := store.New(ctx, params)
	if assert.Nil(t, err) {
		assert.Equal(t, storeID, store.ID())
	}
}

func TestNew_noStoreID(t *testing.T) {

	params := map[string]string{
		store.DriverIdParam: DriverID,
		// no StoreID
	}
	ctx := context.Background()
	store, err := store.New(ctx, params)
	if assert.Nil(t, err) {
		// specific StoreID is unknown: random UUID
		assert.True(t, len(store.ID()) > 0)
	}
}
