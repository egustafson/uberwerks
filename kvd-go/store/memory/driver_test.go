package memory

import (
	"testing"

	"github.com/egustafson/uberwerks/kvd-go/store"
	"github.com/stretchr/testify/assert"
)

func TestDriverID(t *testing.T) {

	driver := new(memStoreDriver)
	assert.Equal(t, DriverID, driver.DriverID())
}

func TestNew(t *testing.T) {

	storeID := "test-store-id"
	params := map[string]string{
		store.DriverIdParam: DriverID,
		store.StoreIdParam:  storeID,
	}
	store, err := store.New(nil, params)
	if assert.Nil(t, err) {
		assert.Equal(t, storeID, store.ID())
	}
}
