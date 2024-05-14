// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package page

import (
	"github.com/webpagine/go-pagine/util"
	"html/template"
	"io"
	"os"
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

	// Parse Go template.
	templateBody, err := os.ReadFile(p.TemplatePath)
	if err != nil {
		return err
	}

	t, err := template.New(p.Path).Parse(string(templateBody))
	if err != nil {
		return err
	}

	// If the page has customized data (encoding in TOML), then add the data it contains.
	if p.DataPath != "" {
		var data map[string]any

		err = util.UnmarshalTOMLFile(p.DataPath, data)
		if err != nil {
			return err
		}

		// Template
		err = t.Execute(w, data) // Yes
		if err != nil {
			return err
		}
	} else {
		// Template
		err = t.Execute(w, nil) // No
		if err != nil {
			return err
		}
	}

	return nil
}
