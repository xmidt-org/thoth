// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"io/fs"

	"github.com/xmidt-org/thoth"
)

type modelLoader struct {
	name   string
	loaded bool
	err    error
	model  thoth.Model
}

// Samples is both a loader and a cache for sample Model data.
type Samples struct {
	Root    fs.FS
	Suffix  string
	samples map[string]interface{}
}

func (s *Samples) Add(name string) {
	if s.samples == nil {
		s.samples = make(map[string]interface{})
	}

	s.samples[name] = struct{}{} // placeholder until actually loaded
}
