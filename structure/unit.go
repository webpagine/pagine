// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"bytes"
	"github.com/webpagine/pagine/v2/collection"
	"github.com/webpagine/pagine/v2/global"
	"github.com/webpagine/pagine/v2/render"
	"github.com/webpagine/pagine/v2/vfs"
	"maps"
	"path/filepath"
	"strings"
	"text/template"
)

type UnitReport struct {
	Error          error
	TemplateErrors []error
}

type UnitManifest struct {
	Units []struct {
		Template string         `toml:"template"`
		Output   string         `toml:"output"`
		Define   map[string]any `toml:"define"`
	} `toml:"unit"`
}

type Unit struct {
	Output   string
	Template string
	Report   UnitReport
}

func (u *Unit) Generate(env *Env, root, dest vfs.DirFS, data MetadataSet, define map[string]any) ([]error, error) {
	var errors collection.Array[error]

	templateName, templateKey := ParseTemplatePair(u.Template)

	t, ok := env.Templates[templateName]
	if !ok {
		return nil, &TemplateUndefinedError{Name: templateName}
	}

	var (
		base, _         = strings.CutPrefix(root.Path, env.Root.Path)
		templateBase, _ = strings.CutPrefix(t.Root.Path, env.Root.Path)

		attr = map[string]any{
			"unitBase":     base,
			"templateBase": templateBase,
		}
	)

	renderFromPath := func(r render.Renderer, pathStr any) string {
		result, err := render.FromPath(r, root, pathStr.(string))
		if err != nil {
			errors.Append(err)
			return ""
		}
		return result
	}

	appliedTemplates := map[string]struct{}{}

	var funcMap map[string]any
	funcMap = map[string]any{
		"add": add,
		"sub": sub,
		"mul": mul,
		"div": div,
		"mod": mod,

		"hasPrefix":  hasPrefix,
		"trimPrefix": trimPrefix,

		"divideSliceByN": divideSliceByN,
		"mapAsSlice":     mapAsSlice,

		"getAttr": func() any { return attr },
		"getEnv":  func() any { return global.EnvAttr },

		"getMetadata": func() any { return data[templateName] },

		"apply": func(pathStr any, data any) any {
			path := filepath.Join(root.Path, pathStr.(string))
			if _, ok := appliedTemplates[path]; ok {
				errors.Append(&RecursiveInvokeError{Templates: nil})
				return nil
			}
			appliedTemplates[path] = struct{}{}
			defer delete(appliedTemplates, path)

			t, err := template.New(filepath.Base(pathStr.(string))).Funcs(funcMap).ParseFS(root, pathStr.(string))
			if err != nil {
				errors.Append(err)
				return ""
			}
			b := bytes.NewBuffer(nil)
			err = t.Execute(b, data)
			if err != nil {
				errors.Append(err)
				return ""
			}
			return b.String()
		},
		"applyFromEnv": func(nameStr any, data any) any {
			path := nameStr.(string)
			if _, ok := appliedTemplates[path]; ok {
				errors.Append(&RecursiveInvokeError{Templates: nil})
				return nil
			}
			appliedTemplates[path] = struct{}{}
			defer delete(appliedTemplates, path)

			split := strings.Split(nameStr.(string), ":")
			if len(split) != 2 {
				errors.Append(&TemplateUndefinedError{Name: nameStr.(string)})
				return nil
			}
			t, ok := env.Templates[split[0]]
			if !ok {
				errors.Append(&TemplateUndefinedError{Name: split[0]})
				return nil
			}
			buf := bytes.NewBuffer(nil)
			err := t.ExecuteTemplate(buf, funcMap, split[1], data)
			if err != nil {
				errors.Append(err)
				return nil
			}
			return buf.String()
		},
		"embed": func(pathStr any) any {
			b, err := root.ReadFile(pathStr.(string))
			if err != nil {
				errors.Append(err)
				return ""
			}
			return string(b)
		},
		"render": func(pathStr any) any {
			r, ok := render.Renderers[filepath.Ext(pathStr.(string))[1:]]
			if !ok {
				errors.Append(&render.NotFoundError{Path: pathStr.(string)})
				return ""
			}
			return renderFromPath(r, pathStr)
		},
		"renderAsciidoc": func(pathStr any) any { return renderFromPath(render.Asciidoc, pathStr) },
		"renderMarkdown": func(pathStr any) any { return renderFromPath(render.Markdown, pathStr) },
	}

	f, err := dest.CreateFile(filepath.Join(base, u.Output))
	if err != nil {
		return nil, err
	}

	dataMap := maps.Clone(data[templateName])
	collection.MergeRawMap(dataMap, define)

	err = t.ExecuteTemplate(f, funcMap, templateKey, dataMap)
	if err != nil {
		return nil, err
	}

	return errors.RawArray, nil
}
