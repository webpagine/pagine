// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

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
