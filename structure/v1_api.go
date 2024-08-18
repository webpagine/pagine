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
		"isServing":    c.Env.IsServing,
		"unitBase":     c.Env.BaseOf(c.Root),
		"templateBase": c.TemplateBase,
	}
	return c.attr
}

func (c *v1Api) Apply(templateName, templateKey string, data any) any {
	return c.Wrap(func() (string, error) {

		var dataMap collection.Map[string, any]

		// Detect recursion.
		if _, ok := c.AppliedTemplates[templateName]; ok {
			return "", &RecursiveInvokeError{Templates: nil}
		}
		c.AppliedTemplates[templateName] = struct{}{}
		defer delete(c.AppliedTemplates, templateName)

		t, err := c.Env.GetTemplateFromAlias(templateName)
		if err != nil {
			return "", err
		}

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
			Errors:           c.Errors,
		})

		return executeTemplate(t, templateKey, funcMap, dataMap.Raw)
	})
}

func (c *v1Api) ApplyFile(std, path string, data map[string]any) string {
	return c.Wrap(func() (_ string, err error) {

		absolutePath := filepath.Join(c.Root.Path, path)

		t, err := c.Env.GetTemplateFile(absolutePath, std)
		if err != nil {
			return "", err
		}

		buf := bytes.NewBuffer(nil)

		err = t.Execute(buf, data)
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	})
}

func (c *v1Api) Exist(path string) bool { return false }
