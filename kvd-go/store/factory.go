package store

import (
	"context"
	"errors"
	"sync"
)

const (
	DriverIdParam = "driver-id"
)

var (
	driverRegistry = make(map[string]StoreDriver)
	driverRegLock  sync.Mutex

	ErrMissingDriverIdParam = errors.New("driver-id parameter missing")
	ErrNoSuchDriver         = errors.New("no such driver")
)

type StoreDriver interface {
	DriverID() string
	New(ctx context.Context, params map[string]string) (Store, error)
}

func New(ctx context.Context, params map[string]string) (Store, error) {
	driverID, ok := params[DriverIdParam]
	if !ok {
		return nil, ErrMissingDriverIdParam
	}
	driverRegLock.Lock()
	defer driverRegLock.Unlock()
	driver, ok := driverRegistry[driverID]
	if !ok {
		return nil, ErrNoSuchDriver
	}
	store, err := driver.New(ctx, params)
	if err == nil {
		registerStore(store)
	}
	return store, err
}

func RegisterDriver(driver StoreDriver) {
	driverRegLock.Lock()
	defer driverRegLock.Unlock()
	driverRegistry[driver.DriverID()] = driver
}
