package store

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

const (
	watchChDepth      = 10
	dispatcherChDepth = 10
)

type WatchDispatcher struct {
	ctx      context.Context
	mutex    sync.Mutex
	inCh     chan *Event
	watchers map[uuid.UUID]*watcher
}

type watcher struct {
	id        uuid.UUID
	ctx       context.Context
	cancel    context.CancelFunc
	lock      sync.Mutex
	ch        chan *Event
	key       string
	range_end string
}

func NewWatchDispatcher(ctx context.Context) *WatchDispatcher {
	wd := &WatchDispatcher{
		ctx:      ctx,
		inCh:     make(chan *Event, dispatcherChDepth),
		watchers: make(map[uuid.UUID]*watcher),
	}

	// goroutine <- finalizer
	go func(wd *WatchDispatcher) {
		<-wd.ctx.Done()
		wd.mutex.Lock()
		defer wd.mutex.Unlock()
		for _, watch := range wd.watchers {
			watch.cancel()
		}
	}(wd)

	// start the dispatcher
	go wd.dispatcher()

	return wd
}

func (wd *WatchDispatcher) dispatcher() {
	for wd.ctx.Err() == nil {
		var event *Event
		// wait for event
		select {
		case event = <-wd.inCh:
		case <-wd.ctx.Done():
			return
		}
		// dispatch event
		func() {
			// holding the mutex ==> all watchers are writable
			wd.mutex.Lock()
			defer wd.mutex.Unlock()
			if wd.ctx.Err() != nil { // mitigate a possible race
				return
			}
			for _, watch := range wd.watchers {
				select {
				case watch.ch <- event:
				case <-watch.ctx.Done():
					continue
				case <-wd.ctx.Done():
					continue
				}
			}
		}()
	}
}

func (wd *WatchDispatcher) NewWatcher(ctx context.Context, key, range_end string) <-chan *Event {
	ctx, cancel := context.WithCancel(ctx)
	watch := &watcher{
		id:        uuid.New(),
		ctx:       ctx,
		cancel:    cancel,
		ch:        make(chan *Event, watchChDepth),
		key:       key,
		range_end: range_end,
	}
	wd.mutex.Lock()
	defer wd.mutex.Unlock()
	wd.watchers[watch.id] = watch

	// goroutine <- finalizer
	go func(watch *watcher, wd *WatchDispatcher) {
		<-watch.ctx.Done()
		wd.mutex.Lock()
		defer wd.mutex.Unlock()
		delete(wd.watchers, watch.id)
		watch.lock.Lock()
		defer watch.lock.Unlock()
		close(watch.ch)
	}(watch, wd)

	return watch.ch
}

func (wd *WatchDispatcher) SendEvent(event *Event) {
	select {
	case wd.inCh <- event:
	case <-wd.ctx.Done():
	}
}
