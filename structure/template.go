// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"bytes"
	"github.com/webpagine/pagine/v2/vfs"
	"io"
	"text/template"
)

type Template struct {
	Root *vfs.DirFS

	CanonicalName string

	Templates map[string]string

	GoTemplate *template.Template
	GetFuncMap GetFuncMap
}

func (t *Template) ExecuteTemplate(wr io.Writer, funcs map[string]any, key string, data any) error {
	name, ok := t.Templates[key]
	if !ok {
		return &TemplateNotFoundError{Template: t, Key: key}
	}

	goTemplate, err := t.GoTemplate.Clone()
	if err != nil {
		return err
	}

	return goTemplate.Funcs(funcs).ExecuteTemplate(wr, name, data)
}
func executeTemplate(t *Template, key string, funcMap, data map[string]any) (string, error) {
	b := bytes.NewBuffer(nil)
	err := t.ExecuteTemplate(b, funcMap, key, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
