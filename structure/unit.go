package structure

import (
	"fmt"
	"github.com/webpagine/go-pagine/collection"
	"github.com/webpagine/go-pagine/render"
	"github.com/webpagine/go-pagine/vfs"
	"maps"
	"path/filepath"
	"strings"
)

type UnitReport struct {
	Error          error
	TemplateErrors []error
}

type UnitManifest struct {
	Unit []struct {
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
	var errors []error

	funcMap := map[string]any{
		"add":            add,
		"divideSliceByN": divideSliceByN,
		"mapAsSlice":     mapAsSlice,

		"embed": func(pathStr any) string {
			b, err := root.ReadFile(pathStr.(string))
			if err != nil {
				errors = append(errors, err)
				return ""
			}
			return string(b)
		},
		"render": func(pathStr any) string {
			r, ok := render.Renderers[filepath.Ext(pathStr.(string))]
			if !ok {
				errors = append(errors, fmt.Errorf("unknown template %q", pathStr))
				return ""
			}
			result, err := render.FromPath(r, root, pathStr.(string))
			if err != nil {
				errors = append(errors, err)
				return ""
			}
			return result
		},
		"renderMarkdown": func(pathStr any) string {
			result, err := render.FromPath(render.Markdown, root, pathStr.(string))
			if err != nil {
				errors = append(errors, err)
				return ""
			}
			return result
		},
	}

	templateName, templateKey := ParseTemplatePair(u.Template)

	if t, ok := env.Templates[templateName]; ok {
		cut, _ := strings.CutPrefix(root.Path, env.Root.Path)
		f, err := dest.CreateFile(filepath.Join(cut, u.Output))
		if err != nil {
			return nil, err
		}

		dataMap := maps.Clone(data[templateName])
		collection.MergeRawMap(dataMap, define)

		err = t.ExecuteTemplate(f, funcMap, templateKey, dataMap)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, &TemplateUndefinedError{Name: templateName}
	}

	return errors, nil
}
