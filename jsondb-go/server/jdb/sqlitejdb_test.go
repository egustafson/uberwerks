package jdb_test

import (
	"testing"

	"github.com/egustafson/uberwerks/jsondb-go/server/jdb"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type sqliteJDBTestSuite struct {
	memoryJDBTestSuite
	dsn string
}

func TestSqliteJDBTestSuite(t *testing.T) {
	suite.Run(t, new(sqliteJDBTestSuite))
}

func (suite *sqliteJDBTestSuite) SetupTest() {
	const dsn = ":memory:"
	var err error
	suite.dsn = dsn
	suite.jdb, err = jdb.InitSqliteJDB(dsn)
	require.Nil(suite.T(), err)
}

func (suite *sqliteJDBTestSuite) TestInitSqliteJDB() {
	suite.NotNil(suite.jdb)
}
