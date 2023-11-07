package wx

type Repository interface {
	Size() int
	Has(k string)
	Get(k string) interface{}
	Put(k string, v interface{})
	Del(k string)
	Enumerate() []string
}
