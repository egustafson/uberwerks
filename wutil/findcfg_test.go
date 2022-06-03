package wutil_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/egustafson/werks/wutil"
)

// ExampleFindConfig demonstrates locating a configuration file for a
// program named 'appctl'
func ExampleFindConfig() {
	_ = wutil.FindConfig("appctl")
}

func TestFindConfig(t *testing.T) {
	//
	// TODO
	//
	return
}

func TestSimpleConfigSearchPath(t *testing.T) {
	basename := "basename"
	userConfigDir, _ := os.UserConfigDir()
	assert.True(t, len(userConfigDir) > 0) // guard against an empty path

	searchPath := wutil.ConfigSearchPath(basename)
	if assert.True(t, len(searchPath) == 2) {
		assert.Equal(t, fmt.Sprintf("%s/%s/%s.yml", userConfigDir, basename, basename), searchPath[0])
		assert.Equal(t, fmt.Sprintf("%s.yml", basename), searchPath[1])
	}
}

func TestProfileConfigSearchPath(t *testing.T) {
	basename := "basename"
	profile := "profile"
	userConfigDir, _ := os.UserConfigDir()
	assert.True(t, len(userConfigDir) > 0) // guard against an empty path

	searchPath := wutil.ConfigSearchPath(basename,
		wutil.WithProfile(profile))
	if assert.True(t, len(searchPath) == 3) {
		assert.Equal(t, fmt.Sprintf("%s/%s/%s/%s.yml", userConfigDir, basename, profile, basename), searchPath[0])
		assert.Equal(t, fmt.Sprintf("%s/%s/%s-%s.yml", userConfigDir, basename, basename, profile), searchPath[1])
		assert.Equal(t, fmt.Sprintf("%s-%s.yml", basename, profile), searchPath[2])
	}
}

func TestProfileAndExtensionConfigSearchPath(t *testing.T) {
	basename := "basename"
	profile := "profile"
	profileExt := "extension"
	userConfigDir, _ := os.UserConfigDir()
	assert.True(t, len(userConfigDir) > 0) // guard against an empty path

	searchPath := wutil.ConfigSearchPath(basename,
		wutil.WithProfile(profile),
		wutil.WithProfileExtension(profileExt))
	if assert.True(t, len(searchPath) == 3) {
		assert.Equal(t, fmt.Sprintf("%s/%s/%s/%s-%s.yml", userConfigDir, basename, profile, basename, profileExt), searchPath[0])
		assert.Equal(t, fmt.Sprintf("%s/%s/%s-%s-%s.yml", userConfigDir, basename, basename, profile, profileExt), searchPath[1])
		assert.Equal(t, fmt.Sprintf("%s-%s-%s.yml", basename, profile, profileExt), searchPath[2])
	}
}

func TestFileExtensionConfigSearchPath(t *testing.T) {
	basename := "basename"
	fileExt := "ini"
	userConfigDir, _ := os.UserConfigDir()
	assert.True(t, len(userConfigDir) > 0) // guard against an empty path

	searchPath := wutil.ConfigSearchPath(basename, wutil.WithFileExtension(fileExt))
	if assert.True(t, len(searchPath) == 2) {
		assert.Equal(t, fmt.Sprintf("%s/%s/%s.%s", userConfigDir, basename, basename, fileExt), searchPath[0])
		assert.Equal(t, fmt.Sprintf("%s.%s", basename, fileExt), searchPath[1])
	}
}

func TestFileBasenameConfigSearchPath(t *testing.T) {
	basename := "basename"
	fileBasename := "alt-filename"
	userConfigDir, _ := os.UserConfigDir()
	assert.True(t, len(userConfigDir) > 0) // guard against an empty path

	searchPath := wutil.ConfigSearchPath(basename, wutil.WithConfigFileBasename(fileBasename))
	if assert.True(t, len(searchPath) == 2) {
		assert.Equal(t, fmt.Sprintf("%s/%s/%s.yml", userConfigDir, basename, fileBasename), searchPath[0])
		assert.Equal(t, fmt.Sprintf("%s.yml", fileBasename), searchPath[1])
	}
}
