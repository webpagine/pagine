// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package util

import (
	"github.com/BurntSushi/toml"
	"github.com/webpagine/go-pagine/vfs"
)

func UnmarshalTOMLFile(root vfs.DirFS, path string, v any) error {
	b, err := root.ReadFile(path)
	if err != nil {
		return err
	}

	return toml.Unmarshal(b, v)
}
