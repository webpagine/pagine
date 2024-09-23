// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/webpagine/pagine/v2/render"
	"os"
	"path/filepath"
)

type v1Render struct {
	*Context
}

func (c *v1Render) FileByMimeType(mimeType, path string) string {
	renderer, err := render.ByMimeType(mimeType)
	panicOnError(err)

	b, err := os.ReadFile(filepath.Join(c.Root.Path, path))
	panicOnError(err)

	result, err := renderer(b)
	panicOnError(err)

	return result
}

func (c *v1Render) FileByExtName(path string) string {
	renderer, err := render.ByExtName(path)
	panicOnError(err)

	b, err := os.ReadFile(filepath.Join(c.Root.Path, path))
	panicOnError(err)

	result, err := renderer(b)
	panicOnError(err)

	return result
}
