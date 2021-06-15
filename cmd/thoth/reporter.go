package main

import (
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

type Result struct {
	Name string
	Err  error
}

type Reporter interface {
	ReportError(templateName string, templateErr error) error
	ReportResults(templateName string, result ...Result) error
}

type ConsoleReporter struct {
	Output  io.Writer
	Verbose bool
}

func (cr ConsoleReporter) output() io.Writer {
	if cr.Output != nil {
		return cr.Output
	}

	return os.Stdout
}

func (cr ConsoleReporter) ReportError(templateName string, templateErr error) (err error) {
	output := cr.output()
	if cr.Verbose || templateErr != nil {
		_, err = fmt.Fprintln(output, templateName)
	}

	if err == nil && templateErr != nil {
		_, err = fmt.Fprintf(output, "%s%-5.5s\t%s\n", indent, ErrorLabel, templateErr)
	}

	return
}

func (cr ConsoleReporter) ReportResults(templateName string, results ...Result) (err error) {
	output := cr.output()
	headerWritten := false

	// just write the header once
	writeHeader := func() {
		if !headerWritten {
			_, err = fmt.Fprintln(output, templateName)
			headerWritten = true
		}
	}

	if cr.Verbose {
		writeHeader()
	}

	for i := 0; err == nil && i < len(results); i++ {
		r := results[i]
		if r.Err != nil {
			writeHeader()
			if err == nil {
				_, err = fmt.Fprintf(output, "%s%-5.5s\t%s\t%s\n", indent, FailLabel, r.Name, r.Err)
			}
		} else if cr.Verbose {
			// the header will have already been written in this case
			_, err = fmt.Fprintf(output, "%s%-5.5s\t%s\n", indent, PassLabel, r.Name)
		}
	}

	return
}
