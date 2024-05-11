package render

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"io"
)

func init() {
	Renderers["md"] = func(r io.Reader, w io.Writer) error {

		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		doc := parser.NewWithExtensions(parser.CommonExtensions | parser.MathJax | parser.NoEmptyLineBeforeBlock).Parse(b)

		renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})

		_, err = w.Write(markdown.Render(doc, renderer))
		if err != nil {
			return err
		}

		return nil
	}
}
