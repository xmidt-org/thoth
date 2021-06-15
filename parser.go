package thoth

import (
	"fmt"
	htemplate "html/template"
	ttemplate "text/template"
)

const (
	MissingKeyInvalid = "invalid"
	MissingKeyDefault = "default"
	MissingKeyZero    = "zero"
	MissingKeyError   = "error"
	DefaultMissingKey = MissingKeyError
)

// InvalidMissingKeyError indicates that an unrecognized value was used in
// configuration for MissingKey.  The text/template and html/template packages
// will panic in this case, whereas this package uses this error.
type InvalidMissingKeyError struct {
	Value string
}

// Error satisfies the error interface.
func (imke *InvalidMissingKeyError) Error() string {
	return fmt.Sprintf("%s is not a valid value for MissingKey", imke.Value)
}

// ParserConfig is the set of configurable options for building a Parser.
type ParserConfig struct {
	// HTML indicates which template package to use.  If this field is false,
	// which is the default, text/template is used by the returned Parser.
	// If this field is true, html/template is used instead.
	HTML bool `json:"html" yaml:"html"`

	// MissingKey is the "missingkey=..." option.  If unset, error is used.
	// If this field is set to an unrecognized value, an error is raised.
	MissingKey string `json:"missingKey" yaml:"missingKey"`

	// LeftDelim is the left delimiter for pipelines.  If unset, the default "{{" is used.
	LeftDelim string `json:"leftDelim" yaml:"leftDelim"`

	// RightDelim is the right delimiter for pipelines.  If unset, the default "}}" is used.
	RightDelim string `json:"rightDelim" yaml:"rightDelim"`

	// FuncMap is the function map for all templates returned by the Parser.
	FuncMap map[string]interface{} `json:"-" yaml:"-"`

	// MediaType is the media type associated with all rendered templates produced
	// by this parser configuration.  If unset, DefaultMediaType is assumed.
	MediaType string `json:"mediaType" yaml:"mediaType"`
}

// templateOptions determines the options to use in prototype templates
func templateOptions(c ParserConfig) (o []string, err error) {
	switch c.MissingKey {
	case "":
		o = append(o, fmt.Sprintf("missingkey=%s", DefaultMissingKey))

	case MissingKeyInvalid:
		fallthrough
	case MissingKeyDefault:
		fallthrough
	case MissingKeyZero:
		fallthrough
	case MissingKeyError:
		o = append(o, fmt.Sprintf("missingkey=%s", c.MissingKey))

	default:
		err = &InvalidMissingKeyError{Value: c.MissingKey}
	}

	return
}

// newPrototype produces the prototype text/template or html/template using
// the given configuration.
func newPrototype(c ParserConfig) (prototype interface{}, err error) {
	var options []string
	options, err = templateOptions(c)
	if err == nil {
		if c.HTML {
			t := htemplate.New("prototype")
			t.Funcs(c.FuncMap)
			t.Delims(c.LeftDelim, c.RightDelim)
			t.Option(options...)
			prototype = t
		} else {
			t := ttemplate.New("prototype")
			t.Funcs(c.FuncMap)
			t.Delims(c.LeftDelim, c.RightDelim)
			t.Option(options...)
			prototype = t
		}
	}

	return
}

// Parser is a template parser.  A ParserConfig defines all the options
// for tailoring a Parser as desired.
type Parser interface {
	// Parse produces a Template from some parsed content.  The name is
	// optional, and can be the empty string.
	//
	// If a common Model was defined in ParserConfig, the returned Template
	// will also implement ModelTemplate.
	//
	// For golang templates, an internal prototype template is first cloned
	// and then used as the containing template for each parsed template.  This
	// ensures that every template is unaffected by global definitions in other templates.
	Parse(name, content string) (Template, error)
}

// NewParser creates a Parser from a set of configuration options.  Templates returned
// by the parser created by this function will also implement MediaTyper, which will
// associated each Template with the media type specified in the config.
func NewParser(c ParserConfig) (Parser, error) {
	prototype, err := newPrototype(c)
	if err != nil {
		return nil, err
	}

	return golangParser{
		prototype: prototype,
		mediaType: c.MediaType,
	}, nil
}

type golangParser struct {
	// prototype is a text/template.Template or html/template.Template which
	// gets cloned to make new templates.
	prototype interface{}

	// mediaType is the media type used when a template doesn't specify one
	mediaType string
}

func (gp golangParser) Parse(name, content string) (t Template, err error) {
	switch pt := gp.prototype.(type) {
	case *ttemplate.Template:
		var raw *ttemplate.Template
		raw, err = pt.Clone()
		if err == nil {
			raw, err = raw.New(name).Parse(content)
			if err == nil {
				t = MediaTemplate(raw, gp.mediaType)
			}
		}

	case *htemplate.Template:
		var raw *htemplate.Template
		raw, err = pt.Clone()
		if err == nil {
			raw, err = raw.New(name).Parse(content)
			if err == nil {
				t = MediaTemplate(raw, gp.mediaType)
			}
		}

	default:
		panic(fmt.Errorf("%T is not a template", t))
	}

	return
}
