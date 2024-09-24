// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/common"
)

func v1GetFuncMap(c *Context) map[string]any {
	funcMap := map[string]any{

		// Engine interaction API.
		"api": common.WrapObject(&v1Api{Context: c}),

		// Filepath processing.
		"path": common.WrapObject(&common.V1Path{}),

		// Rich text format renderer.
		"render": common.WrapObject(&V1Render{Context: c}),

		// v1Strings processing.
		"strings": common.WrapObject(&common.V1Strings{}),

		"This": func() map[string]any { return c.Data },
	}

	// Arithmetic.
	collection.MergeRawMap(funcMap, common.V1Arithmetic)

	// Type casting.
	collection.MergeRawMap(funcMap, common.V1Builtin)

	return funcMap
}

func init() {
	Versions["v1"] = v1GetFuncMap
}
