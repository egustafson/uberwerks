package wx_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/egustafson/werks/wx"
)

const (
	timeout = time.Millisecond
	msg0    = "msg-0"
	msg1    = "msg-1"
	msg2    = "msg-2"
)

func TestNewNotifier(t *testing.T) {
	n := wx.NewNotifier()
	if assert.NotNil(t, n) {
		// send to a notifier w/ zero listeners
		n.Notify("test message")
		n.Close()
	}
}

func TestListener(t *testing.T) {
	n := wx.NewNotifier()

	// send before listen - drops on the floor
	n.Notify(msg0)

	// create a listener
	ch1 := n.Listener(context.Background())
	n.Notify(msg1)

	select {
	case msg := <-ch1:
		assert.Equal(t, msg1, msg)
	case <-time.After(timeout):
		assert.Fail(t, "no message received, expected msg1")
	}

	// create a 2nd listener and observe both receiving
	ch2 := n.Listener(context.Background())
	n.Notify(msg2)

	select {
	case msg := <-ch1:
		assert.Equal(t, msg2, msg)
	case <-time.After(timeout):
		assert.Fail(t, "no message received, expected msg2 on ch1")
	}
	select {
	case msg := <-ch2:
		assert.Equal(t, msg2, msg)
	case <-time.After(timeout):
		assert.Fail(t, "no message received, expected msg2 on ch2")
	}

	n.Close()
	// send after close, no messages delivered, chan's closed
	n.Notify(msg0)
	select {
	case _, ok := <-ch1:
		assert.False(t, ok) // channel closed
	default:
		assert.Fail(t, "expected ch1 to be closed")
	}
	select {
	case _, ok := <-ch2:
		assert.False(t, ok) // channel closed
	default:
		assert.Fail(t, "expected ch2 to be closed")
	}
}

func TestCallback(t *testing.T) {
	n := wx.NewNotifier()
	fixtureChan := make(chan wx.Notice)
	fixture := &callbackFixture{out: fixtureChan}

	n.Callback(context.Background(), fixture.callback)
	n.Notify(msg0)
	n.Notify(msg1)
	n.Notify(msg2) // cause overflow: expect to work

	assert.Equal(t, msg0, <-fixtureChan)
	assert.Equal(t, msg1, <-fixtureChan)
	assert.Equal(t, msg2, <-fixtureChan)

	n.Close()
	n.Notify(msg0)
	select {
	case <-fixtureChan: // still open, nothing to read
		assert.Fail(t, "channel expected to block indefinately, and didn't")
	case <-time.After(timeout):
		// success
	}
}

// --  test fixture for callbacks  --

type callbackFixture struct {
	out chan<- wx.Notice
}

func (cb *callbackFixture) callback(notice wx.Notice) {
	cb.out <- notice
}
