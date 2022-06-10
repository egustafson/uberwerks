package wx

import (
	"context"
	"sync"
	"time"
)

type Notice interface{}

type NotifyCallbackFn func(notice Notice)

type Notifier interface {
	Listener(ctx context.Context) <-chan Notice
	Callback(ctx context.Context, cb NotifyCallbackFn)

	Notify(notice Notice)
	Close()
}

// --  Basic Notifier Impl  ------------------------------------------

type notifyReceiver struct {
	ctx    context.Context
	cancel context.CancelFunc
	ch     chan<- Notice
}

type notifier struct {
	lock      sync.Mutex
	receivers map[*notifyReceiver]struct{}
}

func NewNotifier() Notifier {
	return &notifier{
		receivers: make(map[*notifyReceiver]struct{}),
	}
}

func (n *notifier) Close() {
	n.lock.Lock()
	defer n.lock.Unlock()

	for recv := range n.receivers {
		recv.cancel()
		close(recv.ch)
		delete(n.receivers, recv)
	}
}

func (n *notifier) Notify(notice Notice) {
	const timeout = 50 * time.Microsecond
	n.lock.Lock()
	defer n.lock.Unlock()

	for recv := range n.receivers {
		if recv.ctx.Err() != nil {
			close(recv.ch)
			delete(n.receivers, recv)
			continue
		}
		select {
		case <-recv.ctx.Done():
			close(recv.ch)
			delete(n.receivers, recv)
		case recv.ch <- notice:
			// default: just deliver
		case <-time.After(timeout):
			// spin delivery out to a goroutine; avoid blocking.  Out
			// of order delivery happens at this point
			go func() {
				select {
				case <-recv.ctx.Done():
					return
				case recv.ch <- notice:
					// default: just deliver
				}
			}()
		}
	}
}

func (n *notifier) Listener(ctx context.Context) <-chan Notice {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan Notice, 1)
	recv := &notifyReceiver{
		ctx:    ctx,
		cancel: cancel,
		ch:     ch,
	}
	n.lock.Lock()
	defer n.lock.Unlock()

	n.receivers[recv] = struct{}{}
	return ch
}

func (n *notifier) Callback(ctx context.Context, cb NotifyCallbackFn) {
	listenChan := n.Listener(ctx)

	go func() {
		for notice := range listenChan {
			cb(notice)
		}
	}()
}
