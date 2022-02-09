package wutil_test

import (
	"testing"

	"github.com/egustafson/werks/wutil"
)

// ExampleFindCfg demonstrates locating a configuration file for a
// program named 'appctl'
func ExampleFindCfg() {
	_ = wutil.FindCfg("appctl")
}

func TestFindCfg(t *testing.T) {
	//
	// TODO
	//
	return
}
