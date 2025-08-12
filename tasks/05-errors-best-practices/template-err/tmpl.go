package tmpl

import (
	"errors"
	"html/template"
	"io"
)

var (
	errParseTemplate   = errors.New("can't parse template")
	errExecuteTemplate = errors.New("can't execute template")
)

func ParseAndExecuteTemplate(wr io.Writer, name, text string, data any) error {
	t, err := template.New(name).Parse(text)
	if err != nil {
		return errParseTemplate
	}

	if err := t.Execute(wr, data); err != nil {
		return errExecuteTemplate
	}

	return nil
}
