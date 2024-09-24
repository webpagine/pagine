// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/webpagine/pagine/v2/vfs"
)

type GetFuncMap func(b *Context) map[string]any

type Context struct {
	AppliedTemplates map[string]struct{}
	TemplateBase     string

	Env        *Env
	Root, Dest *vfs.DirFS
	Data       map[string]any
	DataSet    MetadataSet
}

var Versions = map[string]GetFuncMap{}
