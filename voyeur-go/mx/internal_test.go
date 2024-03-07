package mx_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/werks/voyeur-go/mx"
)

func ExampleInternal() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	mx.Internal("demo-internal")
	// Output: [ERROR]   demo-internal [err-file=/home/ericg/voyeur-go/mx/internal_test.go:19][err-func=github.com/werks/voyeur-go/mx_test.ExampleInternal]
}

func TestInternal(t *testing.T) {
	cause := errors.New("stub-causing-error")
	e := mx.Internal("test-internal-error", mx.WithCause(cause), mx.WithMeta("test-meta", "value"))

	assert.Equal(t, "test-internal-error", e.Error())
	assert.True(t, errors.Is(e, cause))
	if ie, ok := e.(*mx.InternalError); assert.True(t, ok) {
		assert.True(t, ie.When.Before(time.Now()))
		value, ok := ie.Metadata["test-meta"]
		if assert.True(t, ok) {
			assert.Equal(t, "value", value)
		}
	}
}

type mockInternalHandler struct{}

func (h *mockInternalHandler) Handle(ie *mx.InternalError) {
	fmt.Println(ie.Message)
}

func ExampleRegisterInternalHandler() {
	mx.RegisterInternalHandler(&mockInternalHandler{})
	mx.Internal("mock-handler-internal-error")
	// Output: mock-handler-internal-error
}
