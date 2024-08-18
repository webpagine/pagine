// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/webpagine/pagine/v2/structure"
	"github.com/webpagine/pagine/v2/vfs"
	"text/template"
)

type TemplateManifest struct {
	Manifest struct {
		Canonical string   `yaml:"canonical"`
		Patterns  []string `yaml:"patterns"`
		Std       string   `yaml:"std"`
	} `yaml:"manifest"`

	Templates []struct {
		Name   string `yaml:"name"`
		Export string `yaml:"export"`
	} `yaml:"templates"`
}

func LoadTemplate(root *vfs.DirFS) (*structure.Template, error) {

	var manifest TemplateManifest

	err := UnmarshalYAMLFile(root, "/manifest.yaml", &manifest)
	if err != nil {
		return nil, err
	}

	exported := map[string]string{}
	for _, t := range manifest.Templates {
		exported[t.Name] = t.Export
	}

	getFuncMap, ok := structure.Versions[manifest.Manifest.Std]
	if !ok {
		return nil, &structure.UndefinedStdError{Std: manifest.Manifest.Std}
	}

	goTemplate, err := template.New(manifest.Manifest.Canonical).Funcs(getFuncMap(nil)).ParseFS(root, manifest.Manifest.Patterns...)
	if err != nil {
		return nil, err
	}

	return &structure.Template{
		Root:          root,
		CanonicalName: manifest.Manifest.Canonical,
		Templates:     exported,
		GoTemplate:    goTemplate,
		GetFuncMap:    getFuncMap,
	}, nil
}
