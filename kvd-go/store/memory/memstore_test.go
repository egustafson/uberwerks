package memory

import (
	"context"
	"testing"
	"time"

	"github.com/egustafson/uberwerks/kvd-go/store"
	"github.com/stretchr/testify/suite"
)

const TestStoreID = "test-store-id"

type TestMemStoreSuite struct {
	suite.Suite
	testCtx    context.Context
	testCancel context.CancelFunc
	store      store.Store
}

func TestMemStore(t *testing.T) {
	suite.Run(t, new(TestMemStoreSuite))
}

func (s *TestMemStoreSuite) SetupTest() {
	s.testCtx, s.testCancel = context.WithCancel(context.Background())
	params := map[string]string{
		store.DriverIdParam: DriverID,
		store.StoreIdParam:  TestStoreID,
	}
	var err error
	s.store, err = store.New(s.testCtx, params)
	s.Require().Nil(err)
}

func (s *TestMemStoreSuite) TearDownTest() {
	s.testCancel()
}

func (s *TestMemStoreSuite) TestMemStore_ID() {
	s.Equal(TestStoreID, s.store.ID())
}

func (s *TestMemStoreSuite) TestMemStore_Manager() {
	mgr := s.store.Manager()
	s.NotNil(mgr)
}

func (s *TestMemStoreSuite) TestMemStore_ctx_canceled() {
	s.testCancel() // cancel the store's context
	ctx := context.Background()

	err := s.store.Put(ctx, "k", "v")
	s.NotNil(err)

	_, err = s.store.KeyRange(ctx, "k0", "k9")
	s.NotNil(err)

	_, err = s.store.GetRange(ctx, "k0", "k9")
	s.NotNil(err)

	_, err = s.store.DelRange(ctx, "k0", "k9")
	s.NotNil(err)

	_, err = s.store.WatchRange(ctx, "k0", "k9")
	s.NotNil(err)
}

func (s *TestMemStoreSuite) TestMemStore_empty() {
	kvs, err := s.store.GetRange(s.testCtx, "!", "~")
	if s.Nil(err) {
		s.True(len(kvs) == 0)
	}
	keys, err := s.store.KeyRange(s.testCtx, "!", "~")
	if s.Nil(err) {
		s.True(len(keys) == 0)
	}
	kvs, err = s.store.DelRange(s.testCtx, "!", "~")
	if s.Nil(err) {
		s.True(len(kvs) == 0)
	}
}

func (s *TestMemStoreSuite) TestMemStore_Put() {
	err := s.store.Put(s.testCtx, "key-1", "value-1")
	s.Nil(err)
}

func (s *TestMemStoreSuite) TestMemStore_PutGetDel_one() {
	err := s.store.Put(s.testCtx, "k1", "v1")
	if s.Nil(err) {

		keys, err := s.store.KeyRange(s.testCtx, "!", "~")
		if s.Nil(err) && s.True(len(keys) == 1) {
			s.Equal("k1", keys[0])
		}
		kvs, err := s.store.GetRange(s.testCtx, "!", "~")
		if s.Nil(err) && s.True(len(kvs) == 1) {
			s.Equal("v1", kvs[0].Val)
		}
		kvsDel, err := s.store.DelRange(s.testCtx, "!", "~")
		if s.Nil(err) && s.True(len(kvsDel) == 1) {
			s.Equal(kvs, kvsDel)
		}
	}
}

func (s *TestMemStoreSuite) TestMemStore_Watch_unused() {
	ctx, cancel := context.WithCancel(context.Background())
	ch, err := s.store.WatchRange(ctx, "!", "~")
	s.Nil(err)

	cancel() // <-- causes ch to be closed
	select {
	case _, ok := <-ch:
		s.False(ok)
	case <-time.After(time.Second):
		s.Fail("expected channel to be closed")
	}
}
