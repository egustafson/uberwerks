package kv

import (
	"context"
	"errors"
)

type Key []byte
type Value []byte

type KeyValue struct {
	K Key
	V Value
}

type EventType int

const (
	PutEvent EventType = iota
	DelEvent
)

type Event struct {
	EventType EventType
	Kv        KeyValue
	PrevKv    KeyValue
}

type KV interface {
	Close()
	Dump() string
	Put(k Key, v Value) error
	Get(k Key) (v Value, err error)
	GetPrefix(prefix Key) (kvs []KeyValue, err error)
	Del(k Key) (err error)
	DelPrefix(prefix Key) (keys []Key, err error)
	Watch(ctx context.Context, k Key) (<-chan []Event, error)
	WatchPrefix(ctx context.Context, prefix Key) (<-chan []Event, error)
}

type NoSuchKeyError interface{ error }

func noSuchKeyError() NoSuchKeyError {
	e := errors.New("no such key")
	return NoSuchKeyError(e)
}

type ClosedError interface{ error }

func closedError() ClosedError {
	e := errors.New("kv store closed")
	return ClosedError(e)
}
