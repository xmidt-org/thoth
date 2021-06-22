package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

const (
	ErrorLabel = "ERROR"
	PassLabel  = "PASS"
	FailLabel  = "FAIL"

	indent = "  "
)

type SampleResult struct {
	Name string
	Err  error
}

type TemplateResult struct {
	Name          string
	Err           error
	SampleResults []SampleResult
}

type Logger interface {
	Debug(message string) error
	Error(message string) error
	Result(TemplateResult) error
}

type ConsoleLogger struct {
	Out     io.Writer
	Err     io.Writer
	Verbose bool

	buffer bytes.Buffer
}

func (cr *ConsoleLogger) out() io.Writer {
	if cr.Out != nil {
		return cr.Out
	}

	return os.Stdout
}

func (cr *ConsoleLogger) err() io.Writer {
	if cr.Err != nil {
		return cr.Err
	}

	return os.Stderr
}

func (cr *ConsoleLogger) Debug(message string) (err error) {
	if !cr.Verbose {
		_, err = fmt.Fprintln(cr.out(), message)
	}

	return
}

func (cr *ConsoleLogger) Error(message string) (err error) {
	_, err = fmt.Fprintln(cr.err(), message)
	return
}

func (cr *ConsoleLogger) Result(tr TemplateResult) (err error) {
	cr.buffer.Reset()

	var (
		headerWritten = false

		// closure that ensures a header line is only written one time.
		// since we may not know if a header will need to be written until
		// iterating over the sample results, we can just blindly call this
		// closure before every output.
		headerOnce = func() {
			if !headerWritten {
				fmt.Fprintln(&cr.buffer, tr.Name)
			}

			headerWritten = true
		}
	)

	if cr.Verbose {
		// always write the header for verbose output,
		// which is just the template's name (relative path)
		headerOnce()
	} else if tr.Err != nil {
		headerOnce()
		fmt.Fprintf(&cr.buffer, "%s%-5.5s\t%s\n", indent, ErrorLabel, tr.Err)
	}

	for _, sr := range tr.SampleResults {
		if sr.Err != nil {
			headerOnce()
			fmt.Fprintf(&cr.buffer, "%s%-5.5s\t%s\t%s\n", indent, FailLabel, sr.Name, sr.Err)
		} else if cr.Verbose {
			headerOnce()
			fmt.Fprintf(&cr.buffer, "%s%-5.5s\t%s\n", indent, PassLabel, sr.Name)
		}
	}

	if cr.buffer.Len() > 0 {
		// we've already terminated each line with a newline
		_, err = fmt.Fprint(cr.out(), cr.buffer.String())
	}

	return
}
