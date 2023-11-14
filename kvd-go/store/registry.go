package store

import "sync"

const StoreIdParam = "store-id"

var (
	storeRegistry = make(map[string]Store)
	storeRegLock  sync.Mutex
)

func Lookup(id string) (Store, bool) {
	storeRegLock.Lock()
	defer storeRegLock.Unlock()
	store, ok := storeRegistry[id]
	return store, ok
}

func registerStore(store Store) {
	storeRegLock.Lock()
	defer storeRegLock.Unlock()
	storeRegistry[store.ID()] = store
}
