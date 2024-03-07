package mx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCause(t *testing.T) {
	cause := errors.New("mock-error")
	o := WithCause(cause)
	eo := new(errorOptions)
	o(eo)
	assert.Equal(t, cause, eo.Cause)
}

func TestWithMeta(t *testing.T) {
	o := WithMeta("test-k", "test-v")
	eo := new(errorOptions)
	o(eo)
	if assert.NotNil(t, eo.Metadata) {
		v, ok := eo.Metadata["test-k"]
		if assert.True(t, ok) {
			assert.Equal(t, "test-v", v)
		}
	}
}
