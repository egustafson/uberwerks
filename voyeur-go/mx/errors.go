package mx

type errorOptions struct {
	Cause    error
	Metadata map[string]string
}

// ErrorOption optional parameter type for error functions.
type ErrorOption func(*errorOptions)

func WithCause(e error) ErrorOption {
	return func(eo *errorOptions) {
		eo.Cause = e
	}
}

func WithMeta(k, v string) ErrorOption {
	return func(eo *errorOptions) {
		if eo.Metadata == nil {
			eo.Metadata = make(map[string]string)
		}
		eo.Metadata[k] = v
	}
}
