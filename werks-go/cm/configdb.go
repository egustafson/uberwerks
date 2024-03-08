package cm

type ConfigDB interface {
	ConfigItem(uuid string) *ConfigItem
}
