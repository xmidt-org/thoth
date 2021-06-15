package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/xmidt-org/thoth"
)

type CLI struct {
	Root     string   `optional default:"." name:"root" short:"R" help:"root directory for file traversal"`
	Verbose  bool     `optional default:"false" name:"verbose" short:"v" help:"verbose output"`
	NoCfg    bool     `optional default:"false" name:"no-cfg" help:"ignore any configuration files"`
	Patterns []string `arg required help:"template patterns"`
}

func parseCommandLine(args []string) (cli CLI, err error) {
	var parser *kong.Kong
	parser, err = kong.New(&cli)
	if err == nil {
		_, err = parser.Parse(args)
	}

	return
}

func run(args []string) (int, error) {
	cli, err := parseCommandLine(args)
	if err != nil {
		return 1, err
	}

	s, err := thoth.NewSelector(
		thoth.SelectorConfig{
			Patterns: cli.Patterns,
		},
	)

	if err != nil {
		return 1, err
	}

	root, err := filepath.Abs(os.ExpandEnv(cli.Root))
	if err != nil {
		return 2, err
	}

	_, _, err = Scanner{
		Root: os.DirFS(root),
		Reporter: ConsoleReporter{
			Output:  os.Stdout,
			Verbose: cli.Verbose,
		},
		Selector: s,
	}.Scan()

	if err != nil {
		return 2, err
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
