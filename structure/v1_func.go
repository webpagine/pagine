// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import "github.com/jellyterra/collection-go"

func v1GetFuncMap(c *Context) map[string]any {
	funcMap := map[string]any{

		// Engine interaction API.
		"api": WrapObject(&v1Api{Context: c}),

		// Filepath processing.
		"path": WrapObject(&v1Path{}),

		// Rich text format renderer.
		"render": WrapObject(&v1Render{Context: c}),

		// v1Strings processing.
		"strings": WrapObject(&v1Strings{}),

		"This": func() map[string]any { return c.Data },
	}

	// Arithmetic.
	collection.MergeRawMap(funcMap, v1Arithmetic)

	// Type casting.
	collection.MergeRawMap(funcMap, v1Cast)

	return funcMap
}

func init() {
	Versions["v1"] = v1GetFuncMap
}
