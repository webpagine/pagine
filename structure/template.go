package structure

import (
	"github.com/webpagine/pagine/util"
	"github.com/webpagine/pagine/vfs"
	"io"
	"text/template"
)

type TemplateManifest struct {
	Manifest struct {
		Canonical string `toml:"canonical"`
	} `toml:"manifest"`

	Templates []struct {
		Src    string `toml:"src"`
		Export string `toml:"export"`
	} `toml:"templates"`
}

type Template struct {
	CanonicalName string

	Templates map[string]string

	GoTemplate *template.Template
}

func (t *Template) ExecuteTemplate(wr io.Writer, funcs map[string]any, key string, data map[string]any) error {
	name, ok := t.Templates[key]
	if !ok {
		return &TemplateNotFoundError{Template: t, Want: key}
	}

	goTemplate, err := t.GoTemplate.Clone()
	if err != nil {
		return err
	}

	return goTemplate.Funcs(funcs).ExecuteTemplate(wr, name, data)
}

func LoadTemplate(root vfs.DirFS) (*Template, error) {

	var manifest TemplateManifest

	err := util.UnmarshalTOMLFile(root, "/manifest.toml", &manifest)
	if err != nil {
		return nil, err
	}

	templateSources := make([]string, len(manifest.Templates))
	templateExported := map[string]string{}

	for i, t := range manifest.Templates {
		templateSources[i] = t.Src
		if t.Export != "" {
			templateExported[t.Export] = t.Src
		}
	}

	goTemplate, err := template.New(manifest.Manifest.Canonical).Funcs(emptyFuncMap).ParseFS(root, templateSources...)
	if err != nil {
		return nil, err
	}

	return &Template{
		CanonicalName: manifest.Manifest.Canonical,
		Templates:     templateExported,
		GoTemplate:    goTemplate,
	}, nil
}

var emptyFuncMap = map[string]any{
	"embed":          func(_ string) string { return "" },
	"render":         func(_ string) string { return "" },
	"renderMarkdown": func(_ string) string { return "" },
}
