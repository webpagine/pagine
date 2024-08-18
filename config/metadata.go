// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/webpagine/pagine/v2/structure"
	"github.com/webpagine/pagine/v2/vfs"
	"os"
)

func LoadMetadata(root *vfs.DirFS) (structure.MetadataSet, error) {
	metadataSet := structure.MetadataSet{}

	err := UnmarshalYAMLFile(root, "/metadata.yaml", &metadataSet)
	switch {
	case err == nil:
	case os.IsNotExist(err):
	default:
		return nil, err
	}

	return metadataSet, nil
}
