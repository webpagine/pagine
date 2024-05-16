// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"io"
)

func init() {
	Renderers[".md"] = func(r io.Reader) ([]byte, error) {

		b, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		doc := parser.NewWithExtensions(parser.CommonExtensions | parser.MathJax | parser.NoEmptyLineBeforeBlock).Parse(b)

		renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})

		return markdown.Render(doc, renderer), nil
	}
}
