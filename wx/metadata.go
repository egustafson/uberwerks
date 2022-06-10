package wx

type Metadata map[string]string

func (md Metadata) Has(k string) bool         { return false }
func (md Metadata) HasExact(k, v string) bool { return false }

func (md Metadata) Find(k string) (string, bool) { return "", false }

func (md Metadata) Set(k, v string) {}

func (md Metadata) Remove(k string)         {}
func (md Metadata) RemoveExact(k, v string) {}

func (md Metadata) MarshalJSON() ([]byte, error) { return []byte{}, nil }

func (md Metadata) UnmarshalJSON(b []byte) error { return nil }
