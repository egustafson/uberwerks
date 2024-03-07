package mx

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Masterminds/log-go"
)

type InternalHandler interface {
	Handle(e *InternalError)
}

var internalErrorHandler InternalHandler = NewBootstrapInternalHandler()

// InternalError holds the details of an internal error
type InternalError struct {
	Message  string
	Cause    error
	When     time.Time
	Metadata map[string]string
}

// Error implements the golang error interface
func (e *InternalError) Error() string {
	return e.Message
}

// Unwrap implements error wrapping
func (e InternalError) Unwrap() error {
	return e.Cause
}

// Internal asserts and returns an internal error
func Internal(msg string, opts ...ErrorOption) error {

	// form the internal error
	errOpts := &errorOptions{
		Metadata: make(map[string]string),
	}
	for _, opt := range opts {
		opt(errOpts)
	}

	e := &InternalError{
		Message:  msg,
		When:     time.Now(),
		Cause:    errOpts.Cause,
		Metadata: errOpts.Metadata,
	}
	// add runtime information
	pc, file, line, ok := runtime.Caller(1) // look up the stack '1' frame
	if ok {
		e.Metadata["err-file"] = fmt.Sprintf("%s:%d", file, line)
		e.Metadata["err-func"] = runtime.FuncForPC(pc).Name()
	} else {
		// this should never happen, but all we can do is record it.
		e.Metadata["err-err"] = "internal-internal-err: problem loading stack backtrace"
	}

	// deliver the internal error
	//
	internalErrorHandler.Handle(e)

	// and return the error
	return e
}

// --
// --  Bootstrap Internal Error Handler  ------------------
// --

type bootstrapInternalHandler struct{}

// static check - *bootstrapInternalHandler isA InternalHandler interface
var _ InternalHandler = (*bootstrapInternalHandler)(nil)

func NewBootstrapInternalHandler() InternalHandler {
	return &bootstrapInternalHandler{}
}

func (h *bootstrapInternalHandler) Handle(e *InternalError) {
	fields := make(log.Fields)
	for k, v := range e.Metadata {
		fields[k] = v
	}
	log.Errorw(e.Message, fields)
}

// RegisterInternalHandler (re)sets the internal error handler.  (Dependency
// Injection)
func RegisterInternalHandler(h InternalHandler) {
	internalErrorHandler = h
}
