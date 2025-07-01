// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package thoth

// SelectorConfig represents a set of templates whose relative paths match
// one or more patterns.
type SelectorConfig struct {
	// Patterns are the globs which must match a template's name in order
	// to use this configured parser.
	Patterns []string `json:"patterns" yaml:"patterns"`

	// Parser is the configuration for parsing templates that match any
	// of the configured patterns.
	Parser ParserConfig `json:"parser" yaml:"parser"`
}

// Selector is a strategy for determining how to parse a template based
// on that template's name.  The name is typically a relative path, which
// may be matched via globbing.
type Selector interface {
	// Select chooses a parser based on the template name.  If no parser
	// was found for the template, this method returns (nil, false).
	Select(name string) (Parser, bool)
}

type matchEntry struct {
	m Matcher
	p Parser
}

type matchSelector struct {
	entries []matchEntry
}

func (ms matchSelector) Select(name string) (p Parser, found bool) {
	for _, e := range ms.entries {
		found = e.m.Match(name)
		if found {
			p = e.p
			break
		}
	}

	return
}

// NewSelector constructs a Selector based on the given configurations.
// If an empty configs is passed, the returned Selector won't match
// any template names.
func NewSelector(configs ...SelectorConfig) (Selector, error) {
	ms := &matchSelector{
		entries: make([]matchEntry, len(configs)),
	}

	for i, c := range configs {
		var err error
		ms.entries[i].m, err = ParsePatterns(c.Patterns...)
		if err == nil {
			ms.entries[i].p, err = NewParser(c.Parser)
		}

		if err != nil {
			return nil, err
		}
	}

	return ms, nil
}
