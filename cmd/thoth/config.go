package main

import (
	"os"
	"path/filepath"

	"github.com/xmidt-org/thoth"
	"gopkg.in/yaml.v3"
)

// ConfigFileName is the name of the thoth configuration file.
const ConfigFileName = ".thoth.yaml"

type Config struct {
	// Verbose controls verbose log output.
	Verbose bool `json:"verbose" yaml:"verbose"`

	// Samples is the set of globs that specify files that are sample models
	// for checking templates.  A sample matches a template if it's file name starts
	// with either the base name or full name of the template.
	Samples []string `json:"samples" yaml:"samples"`

	// Templates associates file patterns with parser configurations.
	Templates []thoth.SelectorConfig `json:"templates" yaml:"templates"`
}

// findConfig starts at an absolute root path and attempts to find ConfigFileName
// by traversing up the directory tree.  If no config file is found, this
// function returns the zero value Config and a nil error.
func findConfig(absoluteRoot string) (c Config, err error) {
	current := absoluteRoot
	for {
		if f, openErr := os.Open(filepath.Join(current, ConfigFileName)); openErr == nil {
			defer f.Close()
			d := yaml.NewDecoder(f)
			err = d.Decode(&c)
			return
		}

		// throw away all trailing separators
		for len(current) > 0 && current[len(current)-1] == os.PathSeparator {
			current = current[0 : len(current)-1]
		}

		if len(current) == 0 {
			break
		}

		current, _ = filepath.Split(current)
	}

	return
}
