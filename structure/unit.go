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
	Unit *Unit

	Error          error
	TemplateErrors []error
}

type Unit struct {
	Output, TemplateName, TemplateKey string

	Define map[string]any

	Report UnitReport
}

func (u *Unit) Generate(env *Env, root, dest *vfs.DirFS, dataSet MetadataSet) ([]error, error) {

	var (
		dataMap collection.Map[string, any]
		errors  collection.SyncVector[error]
	)

	t, err := env.GetTemplateFromAlias(u.TemplateName)
	if err != nil {
		return nil, err
	}

	base := env.BaseOf(root)

	// Inherit.
	dataMap.Raw = maps.Clone(dataSet[t.CanonicalName])

	// Override.
	dataMap.MergeRaw(u.Define)

	// Global template base directory.
	templateBase := env.BaseOf(t.Root)

	funcMap := t.GetFuncMap(&Context{
		AppliedTemplates: map[string]struct{}{},
		TemplateBase:     templateBase,
		Env:              env,
		Root:             root,
		Dest:             dest,
		Data:             dataMap.Raw,
		DataSet:          dataSet,
	})

	f, err := dest.CreateFile(filepath.Join(base, u.Output))
	if err != nil {
		return nil, err
	}

	err = t.ExecuteTemplate(f, funcMap, u.TemplateKey, dataMap.Raw)
	if err != nil {
		return nil, err
	}

	return errors.It.Raw, nil
}
