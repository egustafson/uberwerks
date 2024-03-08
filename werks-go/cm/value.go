package cm

// Value represents an abstract value within the context of ConfigItem and
// ManagedObj attributes where the attribute's name is the key and it's value is
// an implementation of Value.
type Value interface {

	//
	//
	IsSingular() bool

	// Collection type(s)
	//
	IsObject() bool
	IsSequence() bool
}

type ValSingular interface {
	AsString() string
}

type Object interface {
	Keys() []string
	KeysFlattened() []string

	AllKeys() []string
	AllKeysFlattened() []string
}

type Sequence interface {
	Length() int
}
