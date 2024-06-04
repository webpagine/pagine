package structure

import (
	. "github.com/webpagine/go-pagine/collection"
	. "github.com/webpagine/go-pagine/path"
	"github.com/webpagine/go-pagine/util"
	"text/template"
)

type TemplateManifest struct {
	Manifest struct {
		Canonical string `toml:"canonical"`
	} `toml:"manifest"`

	Templates map[string]string `toml:"templates"`
}

type Template struct {
	CanonicalName string

	GoTemplate *template.Template
}

func LoadTemplate(root Path) (*Template, error) {

	var manifest TemplateManifest

	err := util.UnmarshalTOMLFile(root.AbsolutePathOf("manifest.toml"), &manifest)
	if err != nil {
		return nil, err
	}

	t, err := template.New(manifest.Templates["main"]).ParseFiles(ValuesOfRawMap(manifest.Templates)...)
	if err != nil {
		return nil, err
	}

	return &Template{
		CanonicalName: manifest.Manifest.Canonical,
		GoTemplate:    t,
	}, nil
}
