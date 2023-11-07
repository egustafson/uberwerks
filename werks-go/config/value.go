package config

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
