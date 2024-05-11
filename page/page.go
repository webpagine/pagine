package page

import (
	"github.com/webpagine/pagine/util"
	"html/template"
	"io"
)

// Page is configuration of single page.
// index.html.pagine => Config => Template => Renderer (md, rtf) => index.html
type Page struct {

	// Path where the webpage attr file is.
	Path string

	// TemplatePath the pages uses.
	TemplatePath string `toml:"template"`

	// Content sources: content.[string] = string
	// It is designed for different part of content used in template.
	//
	// Example:
	// content = "/content_000.md" => invoke "md" generator, assign the output to the key "content"
	Content map[string]string `toml:"content"`

	// Customized data.
	DataPath string `toml:"data"`
}

func (p *Page) Generate(w io.Writer) error {

	t, err := template.New("").ParseFiles(p.TemplatePath)
	if err != nil {
		return err
	}

	// If has customized data.
	if p.DataPath != "" {
		var data map[string]any

		err = util.UnmarshalTOMLFile(p.DataPath, data)
		if err != nil {
			return err
		}

		err = t.Execute(w, data)
		if err != nil {
			return err
		}
	} else {
		err = t.Execute(w, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
