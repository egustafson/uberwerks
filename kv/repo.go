package kv

import (
	"context"
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

type Repo[T any] interface {
	Get(ctx context.Context, k string) (KV[T], error)
	GetPrefix(ctx context.Context, prefix string) (<-chan KV[T], error)

	Del(ctx context.Context, k string) (count int, err error)
	DelPrefix(ctx context.Context, prefix string) (count int, err error)

	Put(ctx context.Context, k string, v T) error

	Watch(ctx context.Context, k string) (<-chan KV[T], error)
	WatchAll(ctx context.Context, prefix string) (<-chan KV[T], error)
}

//
// --  Remote / Marshallable Repositories ----------------------------
//

type Marshallable interface {
	Marshal() []byte
	Unmarshal([]byte) error
}

type RemoteRepo[T Marshallable] Repo[T]
