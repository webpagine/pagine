package structure

import (
	"fmt"
	"github.com/webpagine/pagine/v2/collection"
	"github.com/webpagine/pagine/v2/render"
	"github.com/webpagine/pagine/v2/vfs"
	"maps"
	"path/filepath"
	"strings"
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
	var errors []error

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
		result, err := render.FromPath(render.Asciidoc, root, pathStr.(string))
		if err != nil {
			errors = append(errors, err)
			return ""
		}
		return result
	}

	funcMap := map[string]any{
		"add": add,
		"sub": sub,
		"mul": mul,
		"div": div,
		"mod": mod,

		"divideSliceByN": divideSliceByN,
		"mapAsSlice":     mapAsSlice,

		"getAttr": func() any { return attr },

		"embed": func(pathStr any) any {
			b, err := root.ReadFile(pathStr.(string))
			if err != nil {
				errors = append(errors, err)
				return ""
			}
			return string(b)
		},
		"render": func(pathStr any) any {
			r, ok := render.Renderers[filepath.Ext(pathStr.(string))[1:]]
			if !ok {
				errors = append(errors, fmt.Errorf("unknown template %q", pathStr))
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

	return errors, nil
}
