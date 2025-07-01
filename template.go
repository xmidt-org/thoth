// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package thoth

import (
	"io"
)

// DefaultMediaType is the MIME type assumed when either an object does not
// implement MediaTyper or returns the empty string from its MediaType method.
const DefaultMediaType = "application/json"

// MediaTyper is an optional interface that associates a MIME type with
// an object.  Typically, this interface is implemented by templates produced
// by this package to indicate what the MIME type of the rendered template is.
type MediaTyper interface {
	// MediaType returns the MIME type associated with this object.  If
	// this method returns the empty string, then DefaultMediaType is assumed.
	MediaType() string
}

// MediaType returns the MIME type associated with the given template.
// If no MIME type is associated with the template, DefaultMediaType is returned.
func MediaType(t Template) (mediaType string) {
	if mt, ok := t.(MediaTyper); ok {
		mediaType = mt.MediaType()
	}

	if len(mediaType) == 0 {
		return DefaultMediaType
	}

	return
}

// Template represents a parsed template.  The golang text and HTML templates
// implement this interface.
type Template interface {
	// Name returns the name of this template supplied at parse time.  This will
	// typically be the relative path of this template.
	Name() string

	// Execute renders this template using the supplied data.  The type of
	// data is unspecified and is dependent on the underlying template.
	Execute(output io.Writer, data interface{}) error
}

type mediaTemplate struct {
	Template
	mediaType string
}

func (mt mediaTemplate) MediaType() string {
	return mt.mediaType
}

// MediaTemplate associates a Template, typically a raw golang template,
// with a MIME type.  The returned Template will also implement MediaTyper.
func MediaTemplate(t Template, mediaType string) Template {
	return mediaTemplate{
		Template:  t,
		mediaType: mediaType,
	}
}
