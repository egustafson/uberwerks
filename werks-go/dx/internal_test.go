package dx_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/egustafson/uberwerks/werks-go/dx"
	"github.com/stretchr/testify/assert"
)

func ExampleInternal() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	dx.Internal("demo-internal")
	// Output: ERROR demo-internal metadata.err-file=/home/ericg/voyeur-go/mx/internal_test.go:19 metadata.err-func=github.com/werks/voyeur-go/mx_test.ExampleInternal
}

func TestInternal(t *testing.T) {
	cause := errors.New("stub-causing-error")
	e := dx.Internal("test-internal-error", dx.WithCause(cause), dx.WithMeta("test-meta", "value"))

	assert.Equal(t, "test-internal-error", e.Error())
	assert.True(t, errors.Is(e, cause))
	if ie, ok := e.(*dx.InternalError); assert.True(t, ok) {
		assert.True(t, ie.When.Before(time.Now()))
		value, ok := ie.Metadata["test-meta"]
		if assert.True(t, ok) {
			assert.Equal(t, "value", value)
		}
	}
}

type mockInternalHandler struct{}

func (h *mockInternalHandler) Handle(ie *dx.InternalError) {
	fmt.Println(ie.Message)
}

func ExampleRegisterInternalHandler() {
	dx.RegisterInternalHandler(&mockInternalHandler{})
	dx.Internal("mock-handler-internal-error")
	// Output: mock-handler-internal-error
}
