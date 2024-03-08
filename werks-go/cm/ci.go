package cm

type ConfigItem interface {
	Key() string // should really be a UUID
	Latest() (version int, val Value)
	Get(version int) Value
}
