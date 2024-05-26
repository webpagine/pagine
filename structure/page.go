// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"bytes"
	"github.com/webpagine/go-pagine/collection"
	. "github.com/webpagine/go-pagine/path"
	"github.com/webpagine/go-pagine/util"
	"io"
	"path/filepath"
	"text/template"
)

// Page is configuration of single page.
// index.html.pagine => Config => Template => Renderer (md, rtf) => index.html
type Page struct {

	// RelativePath where the webpage attr file is.
	RelativePath string

	// Include data from extern.
	Include map[string][]string `toml:"include"`

	// Templates the pages uses.
	Templates map[string]string `toml:"templates"`

	// Contents sources: content.[string] = string
	// It is designed for different part of content used in template.
	//
	// Example:
	// content = "/content_000.md" => invoke "md" generator, assign the output to the key "content"
	Contents map[string]string `toml:"contents"`

	// Customized data.
	Define map[string]map[string]any `toml:"define"`
}

func (p *Page) Generate(root Path, w io.Writer) error {

	mainTemplateRelativePath, ok := p.Templates["main"]
	if !ok {
		return &MissingRequiredFieldError{Field: "templates.main"}
	}

	var (
		templateMap = map[string]string{}
		contentMap  = map[string]any{}
		dataMap     = map[string]map[string]any{}
	)

	for contentKey, contentRelativePath := range p.Contents {
		contentAbsolutePath := root.AbsolutePathOf(contentRelativePath)

		content := NewContent(contentAbsolutePath)
		result, err := content.Generate()
		if err != nil {
			return err
		}

		contentMap[contentKey] = string(result)
	}

	for templateKey := range p.Templates {
		dataMap[templateKey] = map[string]any{}
	}

	for templateKey, list := range p.Include {
		for _, includeRelativePath := range list {
			m := map[string]any{}
			err := util.UnmarshalTOMLFile(root.AbsolutePathOf(includeRelativePath), &m)
			if err != nil {
				return err
			}
			collection.MergeRawMap(dataMap[templateKey], m)
		}
	}

	for templateKey, defMap := range p.Define {
		m := dataMap[templateKey]
		for defKey, defValue := range defMap {
			m[defKey] = defValue
		}
	}

	gen := func(templateKey, templateRelativePath string) error {
		path := root.AbsolutePathOf(templateRelativePath)

		t, err := template.New(filepath.Base(path)).ParseFiles(path)
		if err != nil {
			return err
		}

		b := bytes.NewBuffer(nil)

		err = t.Execute(b, map[string]any{
			"contents":  contentMap,
			"templates": templateMap,
			"data":      dataMap[templateKey],
		})
		if err != nil {
			return err
		}

		templateMap[templateKey] = b.String()

		return nil
	}

	delete(p.Templates, "main")

	for templateKey, templateRelativePath := range p.Templates {
		err := gen(templateKey, templateRelativePath)
		if err != nil {
			return err
		}
	}

	err := gen("main", mainTemplateRelativePath)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(templateMap["main"]))
	if err != nil {
		return err
	}

	return nil
}
