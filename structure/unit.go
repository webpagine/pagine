// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/vfs"
	"maps"
	"path/filepath"
)

type UnitReport struct {
	Error          error
	TemplateErrors []error
}

type UnitManifest struct {
	Units []struct {
		Template string         `yaml:"template"`
		Output   string         `yaml:"output"`
		Define   map[string]any `yaml:"define"`
	} `yaml:"unit"`
}

type Unit struct {
	Output   string
	Template string
	Report   UnitReport
}

func (u *Unit) Generate(env *Env, root, dest *vfs.DirFS, dataSet MetadataSet, define map[string]any) ([]error, error) {

	var dataMap collection.Map[string, any]

	templateName, templateKey := ParseTemplatePair(u.Template)

	// Get matched template from `env`.
	t, ok := env.Templates[templateName]
	if !ok {
		return nil, &TemplateUndefinedError{Name: templateName}
	}

	base := env.BaseOf(root)

	// Inherit.
	dataMap.Raw = maps.Clone(dataSet[templateName])

	// Override.
	dataMap.MergeRaw(define)

	// Global template base directory.
	templateBase := env.BaseOf(t.Root)

	funcMap, errors, err := getFuncMap(
		map[string]struct{}{},
		templateBase,
		env,
		root,
		dest,
		dataMap.Raw,
		dataSet,
	)
	if err != nil {
		return nil, err
	}

	f, err := dest.CreateFile(filepath.Join(base, u.Output))
	if err != nil {
		return nil, err
	}

	err = t.ExecuteTemplate(f, funcMap, templateKey, dataMap.Raw)
	if err != nil {
		return nil, err
	}

	return errors.Raw, nil
}
