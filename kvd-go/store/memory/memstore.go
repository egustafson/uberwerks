package memory

import (
	"context"
	"sync"

	"github.com/egustafson/uberwerks/kvd-go/store"
)

type memStore struct {
	id              string
	ctx             context.Context
	mutex           sync.RWMutex
	revision        uint64
	store           map[string]store.KeyValue
	watchDispatcher *store.WatchDispatcher
}

func newMemStore(ctx context.Context, id string) *memStore {
	st := &memStore{
		id:              id,
		ctx:             ctx,
		revision:        0,
		store:           make(map[string]store.KeyValue),
		watchDispatcher: store.NewWatchDispatcher(ctx),
	}
	return st
}

func (st *memStore) ID() string {
	return st.id
}

func (st *memStore) Put(ctx context.Context, key, value string) error {
	if st.ctx.Err() != nil {
		return st.ctx.Err()
	}
	st.mutex.Lock()
	defer st.mutex.Unlock()
	st.revision += 1
	old_kv, ok := st.store[key]
	if !ok {
		old_kv.Key = key
		// old_kv.Version == 0 => it didn't exist
	}
	kv := store.KeyValue{
		Key:            key,
		Val:            value,
		CreateRevision: old_kv.CreateRevision,
		ModRevision:    st.revision,
		Version:        old_kv.Version + 1,
	}
	if old_kv.CreateRevision == 0 {
		kv.CreateRevision = st.revision
	}
	st.store[key] = kv

	st.watchDispatcher.SendEvent(&store.Event{
		Type:   store.PUT,
		Kv:     kv,
		PrevKv: old_kv,
	})

	return nil
}

func (st *memStore) KeyRange(ctx context.Context, key, range_end string) ([]string, error) {
	if st.ctx.Err() != nil {
		return nil, st.ctx.Err()
	}
	keys := make([]string, 0)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	for k := range st.store {
		if k >= key && k < range_end {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (st *memStore) GetRange(ctx context.Context, key, range_end string) (store.Kvs, error) {
	if st.ctx.Err() != nil {
		return nil, st.ctx.Err()
	}
	kvs := make(store.Kvs, 0)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	for k, v := range st.store {
		if k >= key && k < range_end {
			kvs = append(kvs, v)
		}
	}
	return kvs, nil
}

func (st *memStore) DelRange(ctx context.Context, key, range_end string) (store.Kvs, error) {
	if st.ctx.Err() != nil {
		return nil, st.ctx.Err()
	}
	kvs := make(store.Kvs, 0)
	st.mutex.Lock()
	defer st.mutex.Unlock()
	for k, v := range st.store {
		if k >= key && k < range_end {
			kvs = append(kvs, v)
			delete(st.store, k)
		}
	}
	if len(kvs) > 0 {
		st.revision += 1
	}

	for _, kv := range kvs {
		st.watchDispatcher.SendEvent(&store.Event{
			Type:   store.DEL,
			Kv:     store.KeyValue{Key: kv.Key},
			PrevKv: kv,
		})
	}

	return kvs, nil
}

func (st *memStore) WatchRange(ctx context.Context, key, range_end string) (<-chan *store.Event, error) {
	if st.ctx.Err() != nil {
		return nil, st.ctx.Err()
	}
	return st.watchDispatcher.NewWatcher(ctx, key, range_end), nil
}

func (st *memStore) Manager() any {
	return &MemoryStoreManager{store: st}
}
