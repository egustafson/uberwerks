package wutil

import ()

type CfgLocationOption func(*CfgLocation)

type CfgLocation struct {
	Basename  string
	Profile   string
	Extension string
}

func FindCfg(basename string, opts ...CfgLocationOption) string {

	return ""
}

func ListCfg(basename string, opts ...CfgLocationOption) []string {

	return nil
}
