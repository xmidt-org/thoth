package thoth

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	// Returned by file search functions to indicate that no file was found
	// matching the specified criteria.
	ErrFileNotFound = errors.New("File not found")
)

// UpSearchFile attempts to locate a single, non-directory file by starting
// in a given directory and traversing up to the root.  The names indicate the
// relative file names to search for, and are attempted in order at each level
// of the directory tree.
//
// This function returns the absolute path of any file found, together with its associated
// FileInfo from Stat.  If no file could be found, ErrFileNotFound is returned.
func UpSearchFile(dir string, names ...string) (path string, fi fs.FileInfo, err error) {
	if len(names) == 0 {
		err = ErrFileNotFound
		return
	}

	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return
		}
	}

	for len(dir) > 0 {
		for _, n := range names {
			path = filepath.Join(dir, n)
			fi, err = os.Stat(path)
			if err == nil && !fi.IsDir() {
				return
			}
		}

		for dir[len(dir)-1] == os.PathSeparator {
			dir = dir[0 : len(dir)-1]
		}

		// use Split to traverse up, as it gives a more consistent
		// result for the directory portion
		dir, _ = filepath.Split(dir)
	}

	err = ErrFileNotFound
	return
}
