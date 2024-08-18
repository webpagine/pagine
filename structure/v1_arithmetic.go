// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

var v1Arithmetic = map[string]any{
	"Add": func(aInt, bInt int) int { return aInt + bInt },
	"Sub": func(aInt, bInt int) int { return aInt - bInt },
	"Mul": func(aInt, bInt int) int { return aInt * bInt },
	"Div": func(aInt, bInt int) int { return aInt / bInt },
	"Mod": func(aInt, bInt int) int { return aInt % bInt },
}
