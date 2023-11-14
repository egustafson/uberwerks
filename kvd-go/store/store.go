package store

import "context"

type KeyValue struct {
	Key string
	Val string
	Ver int
}

type Kvs []KeyValue

type EventType int

const (
	PUT EventType = 0
	DEL EventType = 1
)

type Event struct {
	Type   EventType
	Kv     KeyValue
	PrevKv KeyValue
}

type Store interface {
	ID() string
	Put(ctx context.Context, kv KeyValue) error
	KeyRange(ctx context.Context, key, range_end string) ([]string, error)
	GetRange(ctx context.Context, key, range_end string) (Kvs, error)
	DelRange(ctx context.Context, key, range_end string) (Kvs, error)
	WatchRange(ctx context.Context, key, range_end string) (<-chan Event, error)

	Manager() any
}
