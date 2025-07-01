// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package thoth

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	// ErrStopSearch is a sentinel indicating that UpSearch should halt and return.
	ErrStopSearch = errors.New("sentinel error indicating a search should end")
)

// UpSearch walks up a directory applying predicates to each absolute path.
// The given directory is walked up to the root.  If any of the predicates
// return the special value ErrStopSearch, then UpSearch returns immediately
// with a nil error.  Otherwise, any error halts the search and that error is returned.
func UpSearch(dir string, fns ...func(string) error) (err error) {
	if len(fns) == 0 {
		return
	}

	dir, err = filepath.Abs(dir)
	for len(dir) > 0 {
		for _, fn := range fns {
			err = fn(dir)
			if errors.Is(err, ErrStopSearch) {
				err = nil
				return
			} else if err != nil {
				return
			}
		}

		for len(dir) > 0 && dir[len(dir)-1] == os.PathSeparator {
			dir = dir[0 : len(dir)-1]
		}

		// use Split to traverse up, as it gives a more consistent
		// result for the directory portion
		dir, _ = filepath.Split(dir)
	}

	return
}

// FirstFile returns a predicate for UpSearch that matches the first file with
// any of a set of names.  Upon any match, the returned predicate returns ErrStopSearch.
//
// The result pointer must be non-nil and receives the absolute path of the first
// file matched.  Info is optional, and if non-nil it receives the fs.FileInfo associated
// with the result path.
func FirstFile(result *string, info *fs.FileInfo, names ...string) func(string) error {
	return func(dir string) error {
		for _, n := range names {
			path := filepath.Join(dir, n)
			if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
				*result = path
				if info != nil {
					*info = fi
				}

				return ErrStopSearch
			}
		}

		return nil
	}
}
