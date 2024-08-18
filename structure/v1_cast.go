// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

var v1Cast = map[string]any{
	"Any": func(v any) any { return v },
	"Int": func(v any) int { return v.(int) },
	"Str": func(v any) string { return v.(string) },
}
