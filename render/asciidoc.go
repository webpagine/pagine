package render

import (
	"bytes"
	"github.com/bytesparadise/libasciidoc"
	"github.com/bytesparadise/libasciidoc/pkg/configuration"
)

func Asciidoc(content []byte) (string, error) {

	b := bytes.NewBuffer(nil)

	_, err := libasciidoc.Convert(bytes.NewReader(content), b, &configuration.Configuration{})
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
