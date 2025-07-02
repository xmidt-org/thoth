// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package thoth

import (
	"os"

	"github.com/gobwas/glob"
)

// Matcher is a simple strategy for matching values, such as names
// and paths.
type Matcher interface {
	Match(string) bool
}

// Matchers is an aggregate Matcher.  A value will
// match if any of the sequence of Matchers returns true.
type Matchers []Matcher

func (ms Matchers) Match(v string) bool {
	for _, m := range ms {
		if m.Match(v) {
			return true
		}
	}

	return false
}

// ParsePatterns parses a sequence of globs for matching values.  The returned
// Matcher will match values if at least one of the globs matched.  If patterns
// is empty, then the returned Matcher won't match anything.
func ParsePatterns(patterns ...string) (Matcher, error) {
	ms := make(Matchers, len(patterns))
	for _, p := range patterns {
		g, err := glob.Compile(p, os.PathSeparator)
		if err != nil {
			return nil, err
		}

		ms = append(ms, g)
	}

	return ms, nil
}
