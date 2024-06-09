// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func Markdown(content []byte) (string, error) {

	doc := parser.NewWithExtensions(parser.CommonExtensions | parser.MathJax | parser.NoEmptyLineBeforeBlock).Parse(content)

	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})

	return string(markdown.Render(doc, renderer)), nil
}

func init() {
	Renderers[".md"] = Markdown
}
