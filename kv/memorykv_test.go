package kv_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/egustafson/werks/kv"
)

func TestKVTestSuite(t *testing.T) {
	suite.Run(t, new(KVTestSuite))
}

type KVTestSuite struct {
	suite.Suite
	kvdb kv.KV
}

func (s *KVTestSuite) SetupTest() {
	s.kvdb = kv.NewMemoryKV()
}

func (s *KVTestSuite) TearDownTest() {
	s.kvdb.Close()
	s.kvdb = nil
}

func (s *KVTestSuite) TestPutGetDelGet() {
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

func (s *KVTestSuite) TestClose() {
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

func (s *KVTestSuite) TestGetPrefix() {
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

func (s *KVTestSuite) TestDelPrefix() {
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

func (s *KVTestSuite) TestWatch() {
	k := kv.Key("watched-key")
	v := kv.Value("watched-value")
	ctx, cancel := context.WithCancel(context.Background())
	evCh, err := s.kvdb.Watch(ctx, k)
	s.Nil(err)

	err = s.kvdb.Put(k, v)
	s.Nil(err)

	select {
	case events, ok := <-evCh:
		if s.True(ok) {
			s.True(len(events) == 1)
			ev := events[0]
			s.Equal(kv.PutEvent, ev.EventType)
			s.Equal(k, ev.Kv.K)
			s.Equal(v, ev.Kv.V)
			s.Nil(ev.PrevKv.K) // nil key == no previous value
			s.Nil(ev.PrevKv.V)
		}
	case <-time.After(time.Millisecond):
		s.Fail("event never arrived in event channel")
	}

	cancel()
	select {
	case _, ok := <-evCh:
		s.False(ok)
	case <-time.After(time.Millisecond):
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
