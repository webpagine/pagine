// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/structure"
	"github.com/webpagine/pagine/v2/vfs"
	"os"
)

type UnitManifest struct {
	Units []struct {
		Template    string         `yaml:"template"`
		TemplateKey string         `yaml:"template_key"`
		Output      string         `yaml:"output"`
		Define      map[string]any `yaml:"define"`
	} `yaml:"unit"`
}

func LoadUnits(root *vfs.DirFS) ([]*structure.Unit, error) {
	var (
		unitManifest UnitManifest

		units collection.Vector[*structure.Unit]
	)

	err := UnmarshalYAMLFile(root, "unit.yaml", &unitManifest)
	switch {
	case err == nil:
	case os.IsNotExist(err):
	default:
		return nil, err
	}

	for _, unitItem := range unitManifest.Units {
		unit := &structure.Unit{
			Output:       unitItem.Output,
			TemplateName: unitItem.Template,
			TemplateKey:  unitItem.TemplateKey,
			Define:       unitItem.Define,
		}

		units.Push(unit)
	}

	return units.Raw, nil
}
