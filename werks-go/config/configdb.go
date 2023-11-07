package config

type ConfigDB interface {
	ConfigItem(uuid string) *ConfigItem
}
