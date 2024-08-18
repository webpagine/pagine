// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"encoding/json"
	"github.com/go-yaml/yaml"
	"github.com/webpagine/pagine/v2/vfs"
)

func UnmarshalMap(m map[string]any, v any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func UnmarshalYAMLFile(root *vfs.DirFS, path string, v any) error {
	b, err := root.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, v)
}
