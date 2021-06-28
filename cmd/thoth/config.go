package main

import (
	"os"

	"github.com/xmidt-org/thoth"
	"gopkg.in/yaml.v3"
)

// ConfigFileName is the name of the thoth configuration file.
const ConfigFileName = ".thoth.yaml"

type Config struct {
	// Samples is the set of globs that specify files that are sample models
	// for checking templates.  A sample matches a template if it's file name starts
	// with either the base name or full name of the template.
	Samples []string `json:"samples" yaml:"samples"`

	// Templates associates file patterns with parser configurations.
	Templates []thoth.SelectorConfig `json:"templates" yaml:"templates"`
}

// readConfig unmarshals a Config object from the given system file path.
func readConfig(path string) (c Config, err error) {
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		d := yaml.NewDecoder(f)
		err = d.Decode(&c)
	}

	return
}

// findConfig searches for a configuration file beginning at the given
// directory and traversing up the directory tree to the root.  If no
// file is found, this function returns an empty path and a nil error.
// Otherwise, the path to the configuration file is returned along with
// the results of readConfig.
func findConfig(dir string) (path string, c Config, err error) {
	err = thoth.UpSearch(dir, thoth.FirstFile(&path, nil, ConfigFileName))

	if err == nil && len(path) > 0 {
		c, err = readConfig(path)
	}

	return
}
