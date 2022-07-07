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

	v, _ := tune.String("string-key", "default-value")
	assert.Equal(t, "value", *v)

	// 	assert.Equal(t, "value", tune.GetString("string-key"))
	// 	assert.Equal(t, "", tune.GetString("non-existant"))
	// 	assert.Equal(t, "", tune.GetString("int-key"))
}
