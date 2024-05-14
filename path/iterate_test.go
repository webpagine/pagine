// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package path

import (
	"testing"
)

func TestPath_IterateAsRelative(t *testing.T) {
	paths, err := (&Path{"/etc"}).IterateAsRelative()
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		println(path)
	}
}
