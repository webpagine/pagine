// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"bytes"
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/global"
	"github.com/webpagine/pagine/v2/render"
	"github.com/webpagine/pagine/v2/vfs"
	"maps"
	"path/filepath"
	"strings"
	"text/template"
)

func executeTemplate(t *Template, key string, funcMap, data map[string]any) (string, error) {
	b := bytes.NewBuffer(nil)
	err := t.ExecuteTemplate(b, funcMap, key, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func executeTemplateFile(root *vfs.DirFS, name string, funcMap, data map[string]any) (string, error) {
	t, err := template.New(name).Funcs(funcMap).ParseFS(root, name)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(nil)
	err = t.Execute(b, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func getFuncMap(
	appliedTemplates map[string]struct{},
	templateBase string,

	env *Env,
	root, dest *vfs.DirFS,
	data map[string]any,
	dataSet MetadataSet,
) (
	funcMap map[string]any,
	errors *collection.Vector[error],
	_ error,
) {
	errors = &collection.Vector[error]{}

	base := env.BaseOf(root)

	attr := map[string]any{
		"unitBase":     base,
		"templateBase": templateBase,
	}

	wrap := func(f func() (string, error)) string {
		str, err := f()
		if err != nil {
			errors.Push(err)
			return ""
		}

		return str
	}

	renderFromPath := func(r render.Renderer, pathStr any) string {
		return wrap(func() (string, error) {

			result, err := render.FromPath(r, root, pathStr.(string))
			if err != nil {
				return "", err
			}

			return result, nil
		})
	}

	apply := func(pathStr any, data any) any {
		return wrap(func() (string, error) {

			path := filepath.Join(root.Path, pathStr.(string))

			// Detect recursion.
			if _, ok := appliedTemplates[path]; ok {
				return "", &RecursiveInvokeError{Templates: nil}
			}
			appliedTemplates[path] = struct{}{}
			defer delete(appliedTemplates, path)

			// Template file base directory.
			templateBase, _ := strings.CutPrefix(filepath.Base(path), env.Root.Path)

			funcMap, errors2, err := getFuncMap(
				appliedTemplates,
				templateBase,
				env,
				root,
				dest,
				data.(map[string]any),
				dataSet,
			)
			if err != nil {
				return "", err
			}
			defer errors.Push(errors2.Raw...)

			return executeTemplateFile(root, pathStr.(string), funcMap, data.(map[string]any))
		})
	}

	applyFromEnv := func(nameStr any, data any) any {
		return wrap(func() (string, error) {

			var dataMap collection.Map[string, any]

			path := nameStr.(string)

			// Detect recursion.
			if _, ok := appliedTemplates[path]; ok {
				return "", &RecursiveInvokeError{Templates: nil}
			}
			appliedTemplates[path] = struct{}{}
			defer delete(appliedTemplates, path)

			templateName, templateKey := ParseTemplatePair(nameStr.(string))

			// Get matched template from `env`.
			t, ok := env.Templates[templateName]
			if !ok {
				return "", &TemplateUndefinedError{Name: templateName}
			}

			// Inherit.
			dataMap.Raw = maps.Clone(dataSet[templateName])

			// Override.
			dataMap.MergeRaw(data.(map[string]any))

			// Global template base directory.
			templateBase := env.BaseOf(t.Root)

			funcMap, errors2, err := getFuncMap(
				appliedTemplates,
				templateBase,
				env,
				root,
				dest,
				dataMap.Raw,
				dataSet,
			)
			if err != nil {
				return "", err
			}
			defer errors.Push(errors2.Raw...)

			return executeTemplate(t, templateKey, funcMap, dataMap.Raw)
		})
	}

	applyCanonical := func(canonicalStr any, data any) any {
		return wrap(func() (string, error) {
			short, ok := env.CanonicalNames[canonicalStr.(string)]
			if !ok {
				return "", &TemplateUndefinedError{Name: canonicalStr.(string)}
			}

			return applyFromEnv(short, data).(string), nil
		})
	}

	funcMap = map[string]any{

		// Arithmetic.
		"add": add,
		"sub": sub,
		"mul": mul,
		"div": div,
		"mod": mod,

		// String processing.
		"hasPrefix":  hasPrefix,
		"trimPrefix": trimPrefix,

		// Data collection processing.
		"divideSliceByN": divideSliceByN,
		"mapAsSlice":     mapAsSlice,

		// Environment information.
		"getAttr": func() any { return attr },
		"getEnv":  func() any { return global.EnvAttr },

		// Hierarchy information.
		"getMetadata": func() any { return data },

		// Apply the template file located at `pathStr`.
		"apply": apply,

		// Apply the template defined as `nameStr` in `env`.
		"applyFromEnv": applyFromEnv,

		// Apply the template that the canonical name matched.
		"applyCanonical": applyCanonical,

		// Embed the raw file content located at `pathStr`.
		"embed": func(pathStr any) any {
			return wrap(func() (string, error) {

				b, err := root.ReadFile(pathStr.(string))
				if err != nil {
					return "", err
				}

				return string(b), nil
			})
		},

		// Render the file content located at `pathStr` by extension name.
		"render": func(pathStr any) any {
			return wrap(func() (string, error) {

				r, ok := render.Renderers[filepath.Ext(pathStr.(string))[1:]]
				if !ok {
					return "", &render.NotFoundError{Path: pathStr.(string)}
				}

				return renderFromPath(r, pathStr), nil
			})
		},

		// Render the Asciidoc file content located at `pathStr`.
		"renderAsciidoc": func(pathStr any) any { return renderFromPath(render.Asciidoc, pathStr) },

		// Render the Markdown file content located at `pathStr`.
		"renderMarkdown": func(pathStr any) any { return renderFromPath(render.Markdown, pathStr) },
	}

	return
}
