package kv

import (
	"context"

	"github.com/egustafson/werks/wx"
)

func Init() {}

type KV[T any] struct {
	K         string
	V         T
	CreateRev int64
	ModRev    int64
	Version   int64
	err       error
}

type EventType int

const (
	PutEvent EventType = iota
	DelEvent
)

type Event[T any] struct {
	EventType EventType
	Kv        *KV[T]
	PrevKV    *KV[T]
}

type Iterator[T any] interface {
	HasNext() bool
	GetNext() T
}

type WatchResponse[T any] struct {
	Ev  *Event[T]
	Err error
}

type Repo[T any] interface {
	Keys() (Iterator[string], error)
	HasKey(k string) (bool, error)

	Get(k string) (KV[T], error)
	Del(k string) (count int, err error)
	Put(k string, v T) error

	Watch(ctx context.Context, k string) (wx.Notifier[WatchResponse[T]], error)
	WatchPrefix(ctx context.Context, prefix string) (wx.Notifier[WatchResponse[T]], error)
}

//
// --  Remote / Marshallable Repositories ----------------------------
//

type Marshallable interface {
	Marshal() []byte
	Unmarshal([]byte) error
}

type RemoteRepo[T Marshallable] Repo[T]
