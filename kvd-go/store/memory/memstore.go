package memory

import (
	"context"

	"github.com/egustafson/uberwerks/kvd-go/store"
)

type memStore struct {
	id    string
	ctx   context.Context
	store map[string][]byte
}

func (st *memStore) ID() string {
	return st.id
}

func (st *memStore) Put(ctx context.Context, kv store.KeyValue) error {
	return nil
}

func (st *memStore) KeyRange(ctx context.Context, key, range_end string) ([]string, error) {
	return nil, nil
}

func (st *memStore) GetRange(ctx context.Context, key, range_end string) (store.Kvs, error) {
	return nil, nil
}

func (st *memStore) DelRange(ctx context.Context, key, range_end string) (store.Kvs, error) {
	return nil, nil
}

func (st *memStore) WatchRange(ctx context.Context, key, range_end string) (<-chan store.Event, error) {
	return nil, nil
}

func (st *memStore) Manager() any { return nil }
