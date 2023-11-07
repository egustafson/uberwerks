package meta

// Meta is the basic unit of metadata, a string key and value pair.
type Meta struct {
	K string
	V string
}

type Metadata interface {
	// Len returns the total number of Meta items in the set
	Length() int

	// Count returns the number of values for the given key
	Count(key string) int

	// Has returns true if key has one or more values in the set
	Has(key string) bool

	// HasExact returns true if value is present for key
	HasExact(key, value string) bool

	// Find returns all values for key or false if key is not present
	Find(key string) ([]string, bool)

	// All returns a list all key and values
	All() []Meta

	// Set places value as the only value for key.  If key held multiple,
	// previous values then all are replaced by the single value.
	Set(key, value string)

	// Append adds value to the list of values for key.  If key is empty then
	// Append is equavlent to Set()
	Append(key, value string)

	// RemoveAll removes all values keyed by key
	RemoveAll(key string)

	// RemoveExact removes value from key if it was present, otherwise no change
	// occurs.
	RemoveExact(key, value string)

	// TODO:  JSON/YAML (un)marshal
	// MarshalJSON() ([]byte, error)
	// UnmarshalJSON() (b []byte) error
}

type metadata map[string][]string

type MetaOption func(Metadata)

func NewMetadata(opts ...MetaOption) Metadata {
	md := metadata{}
	for _, o := range opts {
		o(md)
	}
	return md
}

func WithMeta(m Meta) MetaOption {
	return func(md Metadata) {
		md.Append(m.K, m.V)
	}
}

func WithMetaList(ml []Meta) MetaOption {
	return func(md Metadata) {
		for _, m := range ml {
			md.Append(m.K, m.V)
		}
	}
}

func WithMetadata(md Metadata) MetaOption {
	return func(newmd Metadata) {
		for _, m := range md.All() {
			newmd.Append(m.K, m.V)
		}
	}
}

// interface Metadata implementation with *struct metadata

func (md metadata) Length() (length int) {
	for _, values := range md {
		for range values {
			length += 1
		}
	}
	return length
}

func (md metadata) Count(key string) (count int) {
	values, ok := md[key]
	if ok {
		for range values {
			count += 1
		}
	}
	return count
}

func (md metadata) Has(key string) bool {
	_, has := md[key]
	return has
}

func (md metadata) HasExact(key, value string) bool {
	values, ok := md[key]
	if ok {
		for _, val := range values {
			if val == value {
				return true
			}
		}
	}
	return false
}

func (md metadata) Find(key string) ([]string, bool) {
	values, ok := md[key]
	return values, ok
}

func (md metadata) All() []Meta {
	all := make([]Meta, 0, len(md))
	for k, values := range md {
		for _, v := range values {
			all = append(all, Meta{k, v})
		}
	}
	return all
}

func (md metadata) Set(key, value string) {
	md[key] = []string{value}
}

func (md metadata) Append(key, value string) {
	values, ok := md[key]
	if ok {
		md[key] = append(values, value)
	} else {
		md[key] = []string{value}
	}
}

func (md metadata) RemoveAll(key string) {
	delete(md, key)
}

func (md metadata) RemoveExact(key, value string) {
	oldValues, ok := md[key]
	if ok {
		newValues := make([]string, 0, len(oldValues))
		for _, v := range oldValues {
			if v != value { // reinsert all but the match(es)
				newValues = append(newValues, v)
			}
		}
		md[key] = newValues
	}
}
