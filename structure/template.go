// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/webpagine/pagine/v2/util"
	"github.com/webpagine/pagine/v2/vfs"
	"io"
	"text/template"
)

type TemplateManifest struct {
	Manifest struct {
		Canonical string   `toml:"canonical"`
		Patterns  []string `toml:"patterns"`
	} `toml:"manifest"`

	Templates []struct {
		Name   string `toml:"name"`
		Export string `toml:"export"`
	} `toml:"templates"`
}

type Template struct {
	Root vfs.DirFS

	CanonicalName string

	Templates map[string]string

	GoTemplate *template.Template
}

func (t *Template) ExecuteTemplate(wr io.Writer, funcs map[string]any, key string, data map[string]any) error {
	name, ok := t.Templates[key]
	if !ok {
		return &TemplateNotFoundError{Template: t, Want: key}
	}

	goTemplate, err := t.GoTemplate.Clone()
	if err != nil {
		return err
	}

	return goTemplate.Funcs(funcs).ExecuteTemplate(wr, name, data)
}

func LoadTemplate(root vfs.DirFS) (*Template, error) {

	var manifest TemplateManifest

	err := util.UnmarshalTOMLFile(root, "/manifest.toml", &manifest)
	if err != nil {
		return nil, err
	}

	exported := map[string]string{}
	for _, t := range manifest.Templates {
		exported[t.Name] = t.Export
	}

	goTemplate, err := template.New(manifest.Manifest.Canonical).Funcs(emptyFuncMap).ParseFS(root, manifest.Manifest.Patterns...)
	if err != nil {
		return nil, err
	}

	return &Template{
		Root:          root,
		CanonicalName: manifest.Manifest.Canonical,
		Templates:     exported,
		GoTemplate:    goTemplate,
	}, nil
}
