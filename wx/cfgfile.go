package wx

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	DefaultFileExtension = "yml"
)

type ConfigPathsOption func(*ConfigPaths)

type ConfigPaths struct {
	Basename         string
	FileBasename     string
	Profile          string
	ProfileExtension string
	FileExtension    string
	UserConfigDir    string
}

func LocateConfigFile(basename string, opts ...ConfigPathsOption) string {

	searchPath := ConfigSearchPath(basename, opts...)

	for _, cpath := range searchPath {
		_, err := os.Stat(cpath)
		if err != nil {
			continue // there's no file at cpath
		}
		f, err := os.Open(cpath)
		if err != nil {
			log.Printf("Could not open for reading: %s: %s", cpath, err)
			continue // the file is unusable, keep trying
		}
		f.Close()
		return cpath
	}
	return "" // empty string ==> no suitable config file located
}

func WithProfile(profile string) ConfigPathsOption {
	return func(cfg *ConfigPaths) {
		cfg.Profile = profile
	}
}

func WithProfileExtension(ext string) ConfigPathsOption {
	return func(cfg *ConfigPaths) {
		cfg.ProfileExtension = ext
	}
}

func WithFileExtension(ext string) ConfigPathsOption {
	return func(cfg *ConfigPaths) {
		cfg.FileExtension = ext
	}
}

func WithConfigFileBasename(fbase string) ConfigPathsOption {
	return func(cfg *ConfigPaths) {
		cfg.FileBasename = fbase
	}
}

func ConfigSearchPath(basename string, opts ...ConfigPathsOption) (searchPath []string) {
	userConfigDir, _ := os.UserConfigDir()
	cp := ConfigPaths{
		Basename:      basename,
		FileBasename:  basename,
		FileExtension: DefaultFileExtension,
		UserConfigDir: userConfigDir,
	}
	for _, o := range opts {
		o(&cp)
	}

	if len(cp.Profile) > 0 {
		searchPath = make([]string, 3) // allocate
		if len(cp.ProfileExtension) > 0 {

			// eg: .config/basename/profile/bfname-pext.ext
			searchPath[0] = filepath.Join(cp.UserConfigDir, cp.Basename, cp.Profile,
				fmt.Sprintf("%s-%s.%s", cp.FileBasename,
					cp.ProfileExtension, cp.FileExtension))

			// eg: .config/basename/bfname-profile-pext.ext
			searchPath[1] = filepath.Join(cp.UserConfigDir, cp.Basename,
				fmt.Sprintf("%s-%s-%s.%s", cp.FileBasename,
					cp.Profile, cp.ProfileExtension, cp.FileExtension))

			// eg: ./bfname-profile-pext.ext
			searchPath[2] = filepath.Join(".",
				fmt.Sprintf("%s-%s-%s.%s", cp.FileBasename,
					cp.Profile, cp.ProfileExtension, cp.FileExtension))
		} else {

			// eg: .config/basename/profile/bfname.ext
			searchPath[0] = filepath.Join(cp.UserConfigDir, cp.Basename, cp.Profile,
				fmt.Sprintf("%s.%s", cp.FileBasename, cp.FileExtension))

			// eg: .config/basename/bfname-profile.ext
			searchPath[1] = filepath.Join(cp.UserConfigDir, cp.Basename,
				fmt.Sprintf("%s-%s.%s", cp.FileBasename, cp.Profile, cp.FileExtension))

			// eg: ./bfname-profile.ext
			searchPath[2] = filepath.Join(".",
				fmt.Sprintf("%s-%s.%s", cp.FileBasename, cp.Profile, cp.FileExtension))
		}
	} else {
		searchPath = make([]string, 2) // allocate

		// eg: .config/basename/bfname.ext
		searchPath[0] = filepath.Join(cp.UserConfigDir, cp.Basename,
			fmt.Sprintf("%s.%s", cp.FileBasename, cp.FileExtension))

		// eg: ./bfname.ext
		searchPath[1] = filepath.Join(".",
			fmt.Sprintf("%s.%s", cp.FileBasename, cp.FileExtension))
	}
	return searchPath
}
