package tune_test

import (
	_ "embed"
	"testing"

	"github.com/egustafson/werks/tune"
	"github.com/stretchr/testify/assert"
)

//go:embed test/test_tunables.yml
var embededTunables string

func TestInit(t *testing.T) {

	err := tune.Init(embededTunables)
	assert.Nil(t, err)
}
