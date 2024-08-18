// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

import (
	"mime"
	"path/filepath"
)

type Renderer func(content []byte) (string, error)

func ByMimeType(mimeType string) (Renderer, error) {
	r, ok := Renderers[mimeType]
	if !ok {
		return nil, &NoRenderFoundError{MimeType: mimeType}
	}

	return r, nil
}

func ByExtName(path string) (Renderer, error) {
	ext := filepath.Ext(path)

	mediaType := mime.TypeByExtension(ext)
	if mediaType == "" {
		return nil, &UnknownExtError{Ext: ext}
	}

	mimeType, _, _ := mime.ParseMediaType(mediaType)

	return ByMimeType(mimeType)
}

var Renderers = map[string]Renderer{}
