// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"bytes"
	. "github.com/webpagine/go-pagine/path"
	"io"
	"path/filepath"
	"text/template"
)

// Page is configuration of single page.
// index.html.pagine => Config => Template => Renderer (md, rtf) => index.html
type Page struct {

	// RelativePath where the webpage attr file is.
	RelativePath string

	// Templates the pages uses.
	Templates map[string]string `toml:"templates"`

	// Contents sources: content.[string] = string
	// It is designed for different part of content used in template.
	//
	// Example:
	// content = "/content_000.md" => invoke "md" generator, assign the output to the key "content"
	Contents map[string]string `toml:"contents"`

	// Customized data.
	Data map[string]any `toml:"data"`
}

func (p *Page) Generate(root Path, w io.Writer) error {

	var contentMap = map[string]any{}

	for contentKey, contentRelativePath := range p.Contents {
		contentAbsolutePath := root.AbsolutePathOf(contentRelativePath)

		content := NewContent(contentAbsolutePath)
		result, err := content.Generate()
		if err != nil {
			return err
		}

		contentMap[contentKey] = string(result)
	}

	for dataKey, dataValue := range p.Data {
		contentMap[dataKey] = dataValue
	}

	var templateMap = map[string]string{}

	contentMap["templates"] = templateMap

	for templateKey, templateRelativePath := range p.Templates {
		path := root.AbsolutePathOf(templateRelativePath)

		t, err := template.New(filepath.Base(path)).ParseFiles(path)
		if err != nil {
			return err
		}

		b := bytes.NewBuffer(nil)

		err = t.Execute(b, contentMap)
		if err != nil {
			return err
		}

		templateMap[templateKey] = b.String()
	}

	_, err := w.Write([]byte(templateMap["main"]))
	if err != nil {
		return err
	}

	return nil
}
