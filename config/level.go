// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/structure"
	"github.com/webpagine/pagine/v2/vfs"
	"path/filepath"
)

func LoadLevel(env *structure.Env, root *vfs.DirFS, parentData structure.MetadataSet) (*structure.Level, error) {
	units, err := LoadUnits(root)
	if err != nil {
		return nil, err
	}

	metadataSet, err := LoadMetadata(root)
	if err != nil {
		return nil, err
	}

	metadataSet, err = env.GenerateMetadataSet(metadataSet)
	if err != nil {
		return nil, err
	}

	return &structure.Level{
		Root:   root,
		Data:   metadataSet.Inherit(parentData),
		Units:  units,
		Levels: nil,
	}, nil
}

func CollectLevels(env *structure.Env, root *vfs.DirFS, parentData structure.MetadataSet) (*structure.Level, error) {

	var levels collection.Vector[*structure.Level]

	entries, err := root.ReadDir("/")
	if err != nil {
		return nil, err
	}

	level, err := LoadLevel(env, root, parentData)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		level, err := CollectLevels(env, vfs.OsDirFS(filepath.Join(root.Path, entry.Name())), level.Data)
		if err != nil {
			return nil, err
		}

		levels.Push(level)
	}

	level.Levels = levels.Raw

	return level, nil
}
