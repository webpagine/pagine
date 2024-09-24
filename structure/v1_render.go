// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/webpagine/pagine/v2/common"
	"github.com/webpagine/pagine/v2/render"
	"os"
	"path/filepath"
)

type V1Render struct {
	*Context
}

func (c *V1Render) FileByMimeType(mimeType, path string) string {
	renderer, err := render.ByMimeType(mimeType)
	common.PanicOnError(err)
	
	b, err := os.ReadFile(filepath.Join(c.Root.Path, path))
	common.PanicOnError(err)

	result, err := renderer(b)
	common.PanicOnError(err)

	return result
}

func (c *V1Render) FileByExtName(path string) string {
	renderer, err := render.ByExtName(path)
	common.PanicOnError(err)

	b, err := os.ReadFile(filepath.Join(c.Root.Path, path))
	common.PanicOnError(err)

	result, err := renderer(b)
	common.PanicOnError(err)

	return result
}
