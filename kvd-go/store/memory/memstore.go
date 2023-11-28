package memory

import (
	"context"
	"sync"

	"github.com/egustafson/uberwerks/kvd-go/store"
	"github.com/google/uuid"
)

type memStore struct {
	id       string
	ctx      context.Context
	mutex    sync.RWMutex
	revision uint64
	store    map[string]store.KeyValue
	watchers map[uuid.UUID]*watcher
}

type watcher struct {
	id        uuid.UUID
	ctx       context.Context
	cancel    context.CancelFunc
	lock      sync.Mutex
	ch        chan store.Event
	key       string
	range_end string
}

func newMemStore(ctx context.Context, id string) *memStore {
	st := &memStore{
		id:       id,
		ctx:      ctx,
		revision: 0,
		store:    make(map[string]store.KeyValue),
		watchers: make(map[uuid.UUID]*watcher),
	}
	go func(st *memStore) {
		<-st.ctx.Done()
		st.mutex.Lock()
		defer st.mutex.Unlock()
		for _, watcher := range st.watchers {
			watcher.cancel()
		}
	}(st)
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
	kv, ok := st.store[key]
	if !ok {
		kv.Key = key
		kv.CreateRevision = st.revision
	}
	kv.Val = value
	kv.Version += 1
	kv.ModRevision = st.revision
	st.store[key] = kv
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
	return kvs, nil
}

func (st *memStore) WatchRange(ctx context.Context, key, range_end string) (<-chan store.Event, error) {
	if st.ctx.Err() != nil {
		return nil, st.ctx.Err()
	}
	ctx, cancel := context.WithCancel(ctx)
	watch := &watcher{
		id:        uuid.New(),
		ctx:       ctx,
		cancel:    cancel,
		ch:        make(chan store.Event, 2),
		key:       key,
		range_end: range_end,
	}
	st.mutex.Lock()
	defer st.mutex.Unlock()
	st.watchers[watch.id] = watch
	go func(watcher *watcher) {
		<-watcher.ctx.Done()
		watcher.lock.Lock()
		defer watcher.lock.Unlock()
		close(watcher.ch)
	}(watch)
	return watch.ch, nil
}

func (st *memStore) Manager() any {
	return &MemoryStoreManager{store: st}
}
