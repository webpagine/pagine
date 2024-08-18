// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/vfs"
)

func WrapObject(v any) func() any { return func() any { return v } }

func WrapMap(funcMap map[string]any) func() map[string]any {
	return func() map[string]any { return funcMap }
}

type GetFuncMap func(b *Context) map[string]any

type Context struct {
	AppliedTemplates map[string]struct{}
	TemplateBase     string

	Env        *Env
	Root, Dest *vfs.DirFS
	Data       map[string]any
	DataSet    MetadataSet

	Errors *collection.SyncVector[error]
}

func (c *Context) Wrap(f func() (string, error)) string {
	str, err := f()
	if err != nil {
		c.Errors.Push(err)
		return ""
	}

	return str
}

var Versions = map[string]GetFuncMap{}
