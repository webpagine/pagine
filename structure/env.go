// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/webpagine/pagine/v2/vfs"
	"path/filepath"
	"regexp"
	"text/template"
)

type Env struct {
	Root *vfs.DirFS

	Templates map[string]*Template

	CanonicalNames map[string]string

	IgnoreGlobs []*regexp.Regexp

	IsServing bool

	DeployTime string

	CachedTemplates map[string]*template.Template
}

func (e *Env) BaseOf(fs *vfs.DirFS) string {
	base, _ := filepath.Rel(e.Root.Path, fs.Path)
	return "/" + base
}

func (e *Env) GenerateMetadataSet(raw MetadataSet) (MetadataSet, error) {
	metadataSet := MetadataSet{}

	for key, data := range raw {
		canonical, ok := e.CanonicalNames[key]
		if ok {
			metadataSet[canonical] = data
		} else {
			metadataSet[key] = data
		}
	}

	return metadataSet, nil
}

func (e *Env) LoadTemplateFile(path, std string) (*template.Template, error) {
	getFuncMap, ok := Versions[std]
	if !ok {
		return nil, &UndefinedStdError{Std: std}
	}

	t, err := template.New(filepath.Base(path)).Funcs(getFuncMap(nil)).ParseFiles(path)
	if err != nil {
		return nil, err
	}

	e.CachedTemplates[path] = t

	return t, nil
}

func (e *Env) GetTemplateFile(path, std string) (*template.Template, error) {
	t, ok := e.CachedTemplates[path]
	if !ok {
		return e.LoadTemplateFile(path, std)
	}
	return t, nil
}

func (e *Env) GetTemplateFromAlias(alias string) (*Template, error) {

	// Is alias.
	if canonical, ok := e.CanonicalNames[alias]; ok {
		return e.Templates[canonical], nil
	}

	// Is canonical.
	if t, ok := e.Templates[alias]; ok {
		return t, nil
	}

	// Not found.
	return nil, &TemplateUndefinedError{Name: alias}
}
