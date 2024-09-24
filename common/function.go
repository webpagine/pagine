// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package common

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func WrapObject(v any) func() any { return func() any { return v } }
