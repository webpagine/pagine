// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	. "github.com/webpagine/go-pagine/path"
	"github.com/webpagine/go-pagine/util"
	"io"
	"os"
	"text/template"
)

// Page is configuration of single page.
// index.html.pagine => Config => Template => Renderer (md, rtf) => index.html
type Page struct {

	// RelativePath where the webpage attr file is.
	RelativePath string

	// TemplatePath the pages uses.
	TemplatePath string `toml:"template"`

	// Contents sources: content.[string] = string
	// It is designed for different part of content used in template.
	//
	// Example:
	// content = "/content_000.md" => invoke "md" generator, assign the output to the key "content"
	Contents map[string]string `toml:"content"`

	// Customized data.
	DataPath string `toml:"data"`
}

func (p *Page) Generate(root Path, w io.Writer) error {

	// Parse Go template.
	templateBody, err := os.ReadFile(root.AbsolutePathOf(p.TemplatePath))
	if err != nil {
		return err
	}

	t, err := template.New(p.RelativePath).Parse(string(templateBody))
	if err != nil {
		return err
	}

	var contentMap = map[string]any{}

	// If the page has customized data (encoding in TOML), then add the data it contains.
	if p.DataPath != "" {
		err = util.UnmarshalTOMLFile(root.AbsolutePathOf(p.DataPath), contentMap)
		if err != nil {
			return err
		}
	}

	for contentKey, contentRelativePath := range p.Contents {
		contentAbsolutePath := root.AbsolutePathOf(contentRelativePath)

		content := NewContent(contentAbsolutePath)
		result, err := content.Generate()
		if err != nil {
			return err
		}

		contentMap[contentKey] = string(result)
	}

	// Template
	err = t.Execute(w, contentMap) // No
	if err != nil {
		return err
	}

	return nil
}
