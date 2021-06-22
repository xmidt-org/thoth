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
	ErrCfgMismatch = errors.New("Cannot specify a config file and the no-cfg option")
)

// CLI represents the thoth command line.
type CLI struct {
	Root      string   `optional default:"." name:"root" short:"R" help:"root directory for file traversal"`
	Verbose   bool     `optional default:"false" name:"verbose" short:"v" help:"verbose output"`
	NoCfg     bool     `optional default:"false" name:"no-cfg" help:"ignore any configuration files"`
	Cfg       string   `optional name:"cfg" help:"explicit configuration file, instead of searching"`
	Samples   []string `optional name:"samples" short:"s" help:"sample patterns"`
	Templates []string `optional name:"templates" short:"t" help:"template patterns"`
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

// loadConfig uses the command line to locate the optional thoth configuration file.
func loadConfig(cli CLI) (c Config, err error) {
	if cli.NoCfg {
		if len(cli.Cfg) > 0 {
			err = ErrCfgMismatch
		}
	} else if len(cli.Cfg) > 0 {
		c, err = readConfig(cli.Cfg)
	} else {
		c, err = findConfig(cli.Root)
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

func newLogger(cli CLI) Logger {
	return &ConsoleLogger{
		Out:     os.Stdout,
		Err:     os.Stderr,
		Verbose: cli.Verbose,
	}
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

	cfg, err := loadConfig(cli)
	if err != nil {
		return ExitBadConfig, err
	}

	selector, err := newSelector(cli, cfg)
	if err != nil {
		return ExitBadConfig, err
	}

	scanner := newScanner(cli, newLogger(cli), selector)
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
