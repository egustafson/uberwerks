package jdb_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/egustafson/uberwerks/jsondb-go/jsondb"
	"github.com/egustafson/uberwerks/jsondb-go/server/jdb"
)

type memoryJDBTestSuite struct {
	suite.Suite
	jdb jdb.JDB
}

func TestMemoryJDBTestSuite(t *testing.T) {
	suite.Run(t, new(memoryJDBTestSuite))
}

func (suite *memoryJDBTestSuite) SetupTest() {
	suite.jdb = jdb.NewMemoryJDB()
	require.NotNil(suite.T(), suite.jdb)
}

func (suite *memoryJDBTestSuite) TestNewMemoryJDB_NewMemoryJDB() {
	suite.NotNil(suite.jdb)
}

func (suite *memoryJDBTestSuite) TestMemoryJDB_List() {
	list, err := suite.jdb.List()
	if suite.Nil(err) {
		suite.Equal(0, len(list))
	}
	suite.jdb.Put(jsondb.JSONObj{jsondb.IDKey: "1"})
	list, err = suite.jdb.List()
	if suite.Nil(err) {
		suite.Equal(1, len(list))
	}
	suite.jdb.Put(jsondb.JSONObj{jsondb.IDKey: "2"})
	list, err = suite.jdb.List()
	if suite.Nil(err) {
		suite.Equal(2, len(list))
	}
}

func (suite *memoryJDBTestSuite) TestMemoryJDB_PutGetDel() {
	jo, err := suite.jdb.Put(jsondb.JSONObj{"key": "value"})
	suite.Nil(err)
	id := jo.ID()
	suite.NotEqual("", id)

	jo, ok, err := suite.jdb.Get(id)
	if suite.True(ok) && suite.Nil(err) {
		suite.Equal("value", jo["key"])
	}

	ok, err = suite.jdb.Del(id)
	suite.Nil(err)
	suite.True(ok)
	_, ok, err = suite.jdb.Get(id)
	suite.Nil(err)
	suite.False(ok)
}

func (suite *memoryJDBTestSuite) TestMemoryJDB_InitMemoryJDB() {
	jdb := jdb.InitMemoryJDB([]map[string]any{
		{"key1": "val1"},
		{"key2": "val2"},
	})
	list, _ := jdb.List()
	suite.Equal(2, len(list))
}

func (suite *memoryJDBTestSuite) TestMemoryJDB_Dump() {
	suite.jdb.Put(jsondb.JSONObj{"key1": "obj2"})
	suite.jdb.Put(jsondb.JSONObj{"key2": "obj2"})

	raw, _ := suite.jdb.Dump()
	suite.True(len(raw) == 2)
}
