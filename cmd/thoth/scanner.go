package main

import (
	"bytes"
	"io/fs"

	"github.com/xmidt-org/thoth"
)

type Scanner struct {
	Root     fs.FS
	Selector thoth.Selector
	Logger   Logger
}

func (s Scanner) Scan() ([]thoth.Template, Samples, error) {
	var (
		buffer    = bytes.NewBuffer(make([]byte, 0, 1024))
		samples   Samples
		templates []thoth.Template
	)

	err := fs.WalkDir(s.Root, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr == nil && !entry.IsDir() {
			if p, found := s.Selector.Select(path); found {
				var t thoth.Template

				f, err := s.Root.Open(path)
				if err == nil {
					defer f.Close()
					buffer.Reset()
					_, err = buffer.ReadFrom(f)
				}

				if err == nil {
					t, err = p.Parse(path, buffer.String())
				}

				s.Logger.Result(TemplateResult{
					Name: path,
					Err:  err,
				})

				if err == nil {
					templates = append(templates, t)
				}
			}
		}

		return nil // always continue
	})

	return templates, samples, err
}
