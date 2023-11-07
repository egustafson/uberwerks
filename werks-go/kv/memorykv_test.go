package kv_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/egustafson/werks/werks-go/kv"
)

func TestMemoryKVTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryKVTestSuite))
}

type MemoryKVTestSuite struct {
	suite.Suite
	kvdb kv.KV
}

func (s *MemoryKVTestSuite) SetupTest() {
	s.kvdb = kv.NewMemoryKV()
}

func (s *MemoryKVTestSuite) TearDownTest() {
	s.kvdb.Close()
	s.kvdb = nil
}

func (s *MemoryKVTestSuite) TestMemoryKV_PutGetDelGet() {
	k := kv.Key([]byte("key"))
	v := kv.Value([]byte("value"))

	err := s.kvdb.Put(k, v)
	s.Nil(err)

	val, err := s.kvdb.Get(k)
	s.Nil(err)
	s.Equal(v, val)

	err = s.kvdb.Del(k)
	s.Nil(err)

	_, err = s.kvdb.Get(k)
	nskErr := kv.NoSuchKeyError(nil)
	s.ErrorAs(err, &nskErr)
}

func (s *MemoryKVTestSuite) TestMemoryKV_Close() {
	s.kvdb.Close()
	closedError := kv.ClosedError(nil)

	s.kvdb.Dump() // indirectly verify it doesn't panic()

	err := s.kvdb.Put(kv.Key(""), kv.Value(""))
	s.ErrorAs(err, &closedError)

	_, err = s.kvdb.Get(kv.Key(""))
	s.ErrorAs(err, &closedError)

	_, err = s.kvdb.GetPrefix(kv.Key(""))
	s.ErrorAs(err, &closedError)

	err = s.kvdb.Del(kv.Key(""))
	s.ErrorAs(err, &closedError)

	_, err = s.kvdb.DelPrefix(kv.Key(""))
	s.ErrorAs(err, &closedError)
}

func createKVList(prefix string, count int) (kvs []kv.KeyValue) {
	for ii := 0; ii < count; ii++ {
		k := fmt.Sprintf("%s-%d", prefix, ii)
		v := fmt.Sprintf("val-%s", k)
		kvs = append(kvs, kv.KeyValue{K: []byte(k), V: []byte(v)})
	}
	return kvs
}

func loadKV(kvs kv.KV, kvlist []kv.KeyValue) {
	for _, kv := range kvlist {
		kvs.Put(kv.K, kv.V)
	}
}

func (s *MemoryKVTestSuite) TestMemoryKV_GetPrefix() {
	p := "key-b"
	loadKV(s.kvdb, createKVList(p, 10))
	loadKV(s.kvdb, createKVList("key-a", 10))
	loadKV(s.kvdb, createKVList("key-c", 10))

	bKeys, err := s.kvdb.GetPrefix(kv.Key([]byte(p)))
	s.Nil(err)
	s.True(len(bKeys) == 10)
	for _, kv := range bKeys {
		s.True(strings.HasPrefix(string(kv.K), p))
	}

}

func (s *MemoryKVTestSuite) TestMemoryKV_DelPrefix() {
	p := "key-b"
	loadKV(s.kvdb, createKVList(p, 10))
	loadKV(s.kvdb, createKVList("key-a", 10))
	loadKV(s.kvdb, createKVList("key-c", 10))

	bkeys, err := s.kvdb.DelPrefix(kv.Key([]byte(p)))
	s.Nil(err)
	s.True(len(bkeys) == 10)
	for _, k := range bkeys {
		s.True(strings.HasPrefix(string(k), p))
	}

	remainingKeys, err := s.kvdb.GetPrefix(kv.Key([]byte("")))
	s.Nil(err)
	for _, kv := range remainingKeys {
		s.False(strings.HasPrefix(string(kv.K), p))
	}
}

func (s *MemoryKVTestSuite) expect(events []kv.Event, evCh <-chan []kv.Event) {
	var received []kv.Event
	var ok bool
	select {
	case received, ok = <-evCh:
		if !ok {
			if len(events) == 0 {
				// success
				return
			} else {
				s.Failf("event channel closed", "expected %d events", len(events))
				return
			}
		}
	case <-time.After(10 * time.Microsecond):
		s.Fail("timeout, event never arrived in event channel")
		return
	}
	eventMap := make(map[string]kv.Event)
	for _, ev := range events {
		switch ev.EventType {
		case kv.PutEvent:
			eventMap[string(ev.Kv.K)] = ev
		case kv.DelEvent:
			eventMap[string(ev.PrevKv.K)] = ev
		}
	}
	if s.True(len(received) == len(events), "expected %s events", len(events)) {
		for _, recv := range received {
			switch recv.EventType {
			case kv.PutEvent:
				if match, ok := eventMap[string(recv.Kv.K)]; ok {
					s.Equal(match.EventType, recv.EventType)
					s.Equal(match.Kv.K, recv.Kv.K)
					s.Equal(match.Kv.V, recv.Kv.V)
				} else {
					s.Failf("unexpected put", "key: %s", string(recv.Kv.K))
				}
			case kv.DelEvent:
				if match, ok := eventMap[string(recv.PrevKv.K)]; ok {
					s.Equal(match.EventType, recv.EventType)
					s.Equal(match.PrevKv.K, recv.PrevKv.K)
				} else {
					s.Failf("unexpected del", "key: %s", string(recv.PrevKv.K))
				}
			}
		}
	}
}

func (s *MemoryKVTestSuite) TestMemoryKV_Watch() {
	k := kv.Key("watched-key")
	v := kv.Value("watched-value")
	ctx, cancel := context.WithCancel(context.Background())
	evCh, err := s.kvdb.Watch(ctx, k)
	s.Nil(err)

	// Test Put
	err = s.kvdb.Put(k, v)
	s.Nil(err)
	s.expect([]kv.Event{{
		EventType: kv.PutEvent,
		Kv:        kv.KeyValue{K: k, V: v},
	}}, evCh)

	// Test Del
	err = s.kvdb.Del(k)
	s.Nil(err)
	s.expect([]kv.Event{{
		EventType: kv.DelEvent,
		PrevKv:    kv.KeyValue{K: k, V: v},
	}}, evCh)

	// Test evCh blocks when no events should be present
	select {
	case <-evCh:
		s.Fail("channel expected to block with no pending events")
	case <-time.After(100 * time.Microsecond):
		// success
	}

	// Test Watch Cancel
	cancel()
	select {
	case _, ok := <-evCh:
		s.False(ok)
	case <-time.After(100 * time.Microsecond):
		s.Fail("channel expected to close, but didn't")
	}
}

func (s *MemoryKVTestSuite) TestMemoryKV_WatchPrefix() {
	key_prefix := "key-"
	keylist := make(map[string]struct{})
	events := make([]kv.Event, 0)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%s%d", key_prefix, i)
		k := kv.Key(key)
		v := kv.Value(fmt.Sprintf("value-%d", i))
		events = append(events, kv.Event{
			EventType: kv.DelEvent,
			PrevKv:    kv.KeyValue{K: k, V: v},
		})
		// take note of the key
		keylist[key] = struct{}{}
		// and, insert the key into the store
		err := s.kvdb.Put(k, v)
		s.Nil(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	evCh, err := s.kvdb.WatchPrefix(ctx, kv.Key(key_prefix))
	s.Nil(err)

	// prefix delete should trigger a list of deleted keys (unknown order)
	_, err = s.kvdb.DelPrefix(kv.Key(key_prefix))
	s.Nil(err)
	s.expect(events, evCh)

	cancel()
	select {
	case _, ok := <-evCh:
		s.False(ok)
	case <-time.After(100 * time.Microsecond):
		s.Fail("channel expected to close, but didn't")
	}
}

func ExampleKV_Dump() {
	kvs := kv.NewMemoryKV()

	kvs.Put(kv.Key("key-2"), kv.Value("value-2"))
	kvs.Put(kv.Key("key-1"), kv.Value("value-1"))
	kvs.Put(kv.Key("key-3"), kv.Value("general value"))

	fmt.Print(kvs.Dump())
	// Output:
	// key-1 value-1
	// key-2 value-2
	// key-3 general value
}
