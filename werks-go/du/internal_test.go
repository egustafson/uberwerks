package du_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/egustafson/uberwerks/werks-go/du"
)

func ExampleInternal() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	du.Internal("demo-internal")
	// Output: [ERROR]   demo-internal [err-file=/home/ericg/voyeur-go/mx/internal_test.go:19][err-func=github.com/werks/voyeur-go/mx_test.ExampleInternal]
}

func TestInternal(t *testing.T) {
	cause := errors.New("stub-causing-error")
	e := du.Internal("test-internal-error", du.WithCause(cause), du.WithMeta("test-meta", "value"))

	assert.Equal(t, "test-internal-error", e.Error())
	assert.True(t, errors.Is(e, cause))
	if ie, ok := e.(*du.InternalError); assert.True(t, ok) {
		assert.True(t, ie.When.Before(time.Now()))
		value, ok := ie.Metadata["test-meta"]
		if assert.True(t, ok) {
			assert.Equal(t, "value", value)
		}
	}
}

type mockInternalHandler struct{}

func (h *mockInternalHandler) Handle(ie *du.InternalError) {
	fmt.Println(ie.Message)
}

func ExampleRegisterInternalHandler() {
	du.RegisterInternalHandler(&mockInternalHandler{})
	du.Internal("mock-handler-internal-error")
	// Output: mock-handler-internal-error
}
