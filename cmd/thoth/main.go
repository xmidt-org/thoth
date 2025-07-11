// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/xmidt-org/thoth"
)

const (
	// ExitBadCommandLine is the process exit code indicating an invalid command line.
	ExitBadCommandLine = iota + 1

	// ExitBadConfig is the process exit code indicating invalid configuration.
	ExitBadConfig

	// ExitScanFailed is the process exit code indicating the the file system scan
	// for templates and/or samples failed.
	ExitScanFailed
)

var (
	// ErrCfgMismatch indicates that the no-cfg and cfg options were both supplied.
	ErrCfgMismatch = errors.New("cannot specify a config file and the no-cfg option")
)

// CLI represents the thoth command line.
type CLI struct {
	Root      string   `optional:"true" default:"." name:"root" short:"R" help:"root directory for file traversal"`
	Verbose   bool     `optional:"true" default:"false" name:"verbose" short:"v" help:"verbose output"`
	NoCfg     bool     `optional:"true" default:"false" name:"no-cfg" help:"ignore any configuration files"`
	Cfg       string   `optional:"true" name:"cfg" help:"explicit configuration file, instead of searching"`
	Samples   []string `optional:"true" name:"samples" short:"s" help:"sample patterns"`
	Templates []string `optional:"true" name:"templates" short:"t" help:"template patterns"`
}

// parseCommandLine uses kong to parse the given arguments and return the CLI instance.
func parseCommandLine(args []string) (cli CLI, err error) {
	var parser *kong.Kong
	parser, err = kong.New(&cli)
	if err == nil {
		_, err = parser.Parse(args)
	}

	if err == nil {
		cli.Root, err = filepath.Abs(os.ExpandEnv(cli.Root))
	}

	return
}

func newLogger(cli CLI) Logger {
	return &ConsoleLogger{
		Out:     os.Stdout,
		Err:     os.Stderr,
		Verbose: cli.Verbose,
	}
}

// loadConfig uses the command line to locate the optional thoth configuration file.
func loadConfig(cli CLI, l Logger) (c Config, err error) {
	switch {
	case cli.NoCfg:
		if len(cli.Cfg) > 0 {
			err = ErrCfgMismatch
		}

	case len(cli.Cfg) > 0:
		c, err = readConfig(cli.Cfg)
		if err != nil {
			err = fmt.Errorf("unable to read configuration file [%s]: %w", cli.Cfg, err)
		} else {
			l.Debugf("using config file %s", cli.Cfg)
		}

	default:
		var path string
		path, c, err = findConfig(cli.Root)
		if err != nil {
			if len(path) > 0 {
				err = fmt.Errorf("unable to read configuration file [%s]: %w", path, err)
			} else {
				err = fmt.Errorf("unable to search for configuration file: %w", err)
			}
		} else if len(path) > 0 {
			l.Debugf("found config file %s", path)
		}
	}

	return
}

func newSelector(cli CLI, cfg Config) (thoth.Selector, error) {
	var scfgs []thoth.SelectorConfig
	if len(cli.Templates) > 0 {
		// any template globs from the command-line are given
		// the default parser settings
		scfgs = append(scfgs,
			thoth.SelectorConfig{
				Patterns: cli.Templates,
			},
		)
	}

	scfgs = append(scfgs, cfg.Templates...)
	return thoth.NewSelector(scfgs...)
}

func newScanner(cli CLI, r Logger, s thoth.Selector) Scanner {
	return Scanner{
		Root:     os.DirFS(cli.Root),
		Logger:   r,
		Selector: s,
	}
}

func run(args []string) (int, error) {
	cli, err := parseCommandLine(args)
	if err != nil {
		return ExitBadCommandLine, err
	}

	l := newLogger(cli)
	cfg, err := loadConfig(cli, l)
	if err != nil {
		return ExitBadConfig, err
	}

	selector, err := newSelector(cli, cfg)
	if err != nil {
		return ExitBadConfig, err
	}

	scanner := newScanner(cli, l, selector)
	_, _, err = scanner.Scan()
	if err != nil {
		return ExitScanFailed, err
	}

	return 0, nil
}

func main() {
	exit, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(exit)
}
