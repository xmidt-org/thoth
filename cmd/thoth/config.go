package main

import (
	"errors"
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

func readConfig(path string) (c Config, err error) {
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		d := yaml.NewDecoder(f)
		err = d.Decode(&c)
	}

	return
}

func findConfig(dir string) (c Config, err error) {
	path, _, searchErr := thoth.UpSearchFile(dir)
	if errors.Is(searchErr, thoth.ErrFileNotFound) {
		return
	} else if searchErr != nil {
		err = searchErr
		return
	}

	c, err = readConfig(path)
	return
}
