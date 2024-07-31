// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package util

import "github.com/webpagine/pagine/v2/vfs"
import "github.com/go-yaml/yaml"

func UnmarshalYAMLFile(root *vfs.DirFS, path string, data interface{}) error {
	b, err := root.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, data)
}
