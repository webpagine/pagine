// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"bytes"
	"github.com/jellyterra/collection-go"
	"maps"
	"path/filepath"
)

type v1Api struct {
	*Context

	attr map[string]any
}

func (c *v1Api) Attr() map[string]any {
	if c.attr != nil {
		return c.attr
	}
	c.attr = map[string]any{
		"deployTime":   c.Env.DeployTime,
		"isServing":    c.Env.IsServing,
		"unitBase":     c.Env.BaseOf(c.Root),
		"templateBase": c.TemplateBase,
	}
	return c.attr
}

func (c *v1Api) Apply(templateName, templateKey string, data any) any {
	var dataMap collection.Map[string, any]

	// Detect recursion.
	if _, ok := c.AppliedTemplates[templateName]; ok {
		panic(&RecursiveInvokeError{Templates: nil})
	}
	c.AppliedTemplates[templateName] = struct{}{}
	defer delete(c.AppliedTemplates, templateName)

	t, err := c.Env.GetTemplateFromAlias(templateName)
	panicOnError(err)

	// Inherit.
	dataMap.Raw = maps.Clone(c.DataSet[t.CanonicalName])

	// Override.
	dataMap.MergeRaw(data.(map[string]any))

	// Global template base directory.
	templateBase := c.Env.BaseOf(t.Root)

	funcMap := t.GetFuncMap(&Context{
		AppliedTemplates: c.AppliedTemplates,
		TemplateBase:     templateBase,
		Env:              c.Env,
		Root:             c.Root,
		Dest:             c.Dest,
		Data:             dataMap.Raw,
		DataSet:          nil,
	})

	result, err := executeTemplate(t, templateKey, funcMap, dataMap.Raw)
	panicOnError(err)

	return result
}

func (c *v1Api) ApplyFile(std, path string, data map[string]any) string {
	getFuncMap, ok := Versions[std]
	if !ok {
		panic(&UndefinedStdError{Std: std})
	}

	funcMap := getFuncMap(&Context{
		AppliedTemplates: c.AppliedTemplates,
		TemplateBase:     filepath.Base(path),
		Env:              c.Env,
		Root:             c.Root,
		Dest:             c.Dest,
		Data:             data,
		DataSet:          nil,
	})

	absolutePath := filepath.Join(c.Root.Path, path)

	t, err := c.Env.GetTemplateFile(absolutePath, std)
	panicOnError(err)

	buf := bytes.NewBuffer(nil)

	err = t.Funcs(funcMap).Execute(buf, data)
	panicOnError(err)

	return buf.String()
}

func (c *v1Api) Exist(path string) bool { return false }
