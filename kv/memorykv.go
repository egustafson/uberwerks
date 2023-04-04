package kv

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type memoryKV struct {
	kv       map[string]Value
	watchers map[chan []Event]context.CancelFunc
	lock     sync.RWMutex
	alive    context.Context
	cancel   context.CancelFunc
}

var _ KV = (*memoryKV)(nil)

func NewMemoryKV() KV {
	ctx, cancel := context.WithCancel(context.Background())
	return &memoryKV{
		kv:       make(map[string]Value),
		watchers: make(map[chan []Event]context.CancelFunc),
		alive:    ctx,
		cancel:   cancel,
	}
}

func (mkv *memoryKV) Close() {
	mkv.lock.Lock()
	defer mkv.lock.Unlock()
	mkv.cancel()
	mkv.kv = nil
}

type sortedkeys []string

func (k sortedkeys) Len() int           { return len(k) }
func (k sortedkeys) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k sortedkeys) Less(i, j int) bool { return string(k[i]) < string(k[j]) }

func (mkv *memoryKV) Dump() string {
	mkv.lock.RLock()
	defer mkv.lock.RUnlock()
	// Dump ignores the alive context; it's a debug function
	//
	keys := make([]string, 0, len(mkv.kv))
	for k := range mkv.kv {
		keys = append(keys, k)
	}
	sort.Sort(sortedkeys(keys))

	var strBuf bytes.Buffer
	for _, k := range keys {
		strBuf.WriteString(fmt.Sprintf("%s %s\n", k, string(mkv.kv[k])))
	}
	return strBuf.String()
}

func (mkv *memoryKV) Put(k Key, v Value) error {
	mkv.lock.Lock()
	defer mkv.lock.Unlock()
	if mkv.alive.Err() != nil {
		return closedError()
	}
	var prevKv KeyValue
	prevV, ok := mkv.kv[string(k)]
	if ok {
		prevKv = KeyValue{K: k, V: prevV}
	}
	mkv.kv[string(k)] = v
	ev := Event{
		EventType: PutEvent,
		Kv:        KeyValue{K: k, V: v},
		PrevKv:    prevKv,
	}
	mkv.sendEvents([]Event{ev})
	return nil
}

func (mkv *memoryKV) Get(k Key) (v Value, err error) {
	mkv.lock.RLock()
	defer mkv.lock.RUnlock()
	if mkv.alive.Err() != nil {
		return nil, closedError()
	}
	var ok bool
	v, ok = mkv.kv[string(k)]
	if !ok {
		return nil, noSuchKeyError()
	}
	return v, nil
}

func (mkv *memoryKV) GetPrefix(k Key) (kvs []KeyValue, err error) {
	mkv.lock.RLock()
	defer mkv.lock.RUnlock()
	if mkv.alive.Err() != nil {
		return nil, closedError()
	}
	prefix := string(k)
	kvs = make([]KeyValue, 0)
	for k, v := range mkv.kv {
		if strings.HasPrefix(k, prefix) {
			kvs = append(kvs, KeyValue{K: []byte(k), V: v})
		}
	}
	return kvs, nil
}

func (mkv *memoryKV) Del(k Key) (err error) {
	mkv.lock.Lock()
	defer mkv.lock.Unlock()
	if mkv.alive.Err() != nil {
		return closedError()
	}
	prevV, ok := mkv.kv[string(k)]
	if !ok {
		return noSuchKeyError()
	}
	delete(mkv.kv, string(k))
	ev := Event{
		EventType: DelEvent,
		PrevKv:    KeyValue{K: k, V: prevV},
	}
	mkv.sendEvents([]Event{ev})
	return nil
}

func (mkv *memoryKV) DelPrefix(k Key) (keys []Key, err error) {
	mkv.lock.Lock()
	defer mkv.lock.Unlock()
	if mkv.alive.Err() != nil {
		return nil, closedError()
	}
	prefix := string(k)
	keys = make([]Key, 0)
	events := make([]Event, 0)
	for k, v := range mkv.kv {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, []byte(k))
			delete(mkv.kv, k)
			ev := Event{
				EventType: DelEvent,
				PrevKv:    KeyValue{K: Key(k), V: Value(v)},
			}
			events = append(events, ev)
		}
	}
	mkv.sendEvents(events)
	return keys, nil
}

// --  Watchers  ------------------------------------------

// sendEvents will deliver 'events' to all the watchers in channels.
// ASSUMPTION:  mkv.lock is write locked by the caller.
func (mkv *memoryKV) sendEvents(events []Event) {
	for ch := range mkv.watchers {
		select {
		case ch <- events:
		case <-time.After(10 * time.Microsecond):
			// don't block, if ch isn't available,
			// spin this out into a goroutine.
			go func(ch chan []Event, events []Event) {
				ch <- events
			}(ch, events)
		}
	}
}

type matcher func(string) bool

func keyMatcherer(key string) matcher {
	return func(pat string) bool {
		return pat == key
	}
}

func prefixMatcherer(prefix string) matcher {
	return func(pat string) bool {
		return strings.HasPrefix(pat, prefix)
	}
}

func (mkv *memoryKV) watchRelay(ctx context.Context, in, out chan []Event, m matcher) {
	defer close(out)
	for { // ever, (or until a context is canceled or a channel closed)
		select {
		case <-ctx.Done():
			mkv.lock.Lock()
			defer mkv.lock.Unlock()
			close(in)
			delete(mkv.watchers, in)
			return // we're done
		case e := <-in:
			out <- e // relay
		}
	}
}

func (mkv *memoryKV) makeWatch(ctx context.Context, m matcher) (<-chan []Event, error) {
	mkv.lock.Lock()
	defer mkv.lock.Unlock()
	if mkv.alive.Err() != nil {
		return nil, closedError()
	}

	watchCh := make(chan []Event)
	listenCh := make(chan []Event)
	ctx, cancel := context.WithCancel(ctx)
	mkv.watchers[listenCh] = cancel
	go mkv.watchRelay(ctx, listenCh, watchCh, m)

	return watchCh, nil
}

func (mkv *memoryKV) Watch(ctx context.Context, k Key) (<-chan []Event, error) {
	return mkv.makeWatch(ctx, keyMatcherer(string(k)))
}

func (mkv *memoryKV) WatchPrefix(ctx context.Context, k Key) (<-chan []Event, error) {
	return mkv.makeWatch(ctx, prefixMatcherer(string(k)))
}
