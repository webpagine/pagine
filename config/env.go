// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/webpagine/pagine/v2/structure"
	"github.com/webpagine/pagine/v2/vfs"
	"regexp"
	"text/template"
)

type EnvManifest struct {
	Use map[string]string `yaml:"use"`

	Ignore []string `yaml:"ignore"`
}

func LoadEnv(root *vfs.DirFS) (*structure.Env, error) {
	var (
		env = structure.Env{
			Root:            root,
			Templates:       map[string]*structure.Template{},
			CanonicalNames:  map[string]string{},
			CachedTemplates: map[string]*template.Template{},
		}
		manifest EnvManifest
	)

	err := UnmarshalYAMLFile(root, "/env.yaml", &manifest)
	if err != nil {
		return nil, err
	}

	env.IgnoreGlobs = make([]*regexp.Regexp, len(manifest.Ignore))
	for i, globForm := range manifest.Ignore {
		glob, err := regexp.Compile(globForm)
		if err != nil {
			return nil, err
		}
		env.IgnoreGlobs[i] = glob
	}

	for templateAlias, templatePath := range manifest.Use {
		sub, err := root.Chroot(templatePath)
		if err != nil {
			return nil, err
		}

		t, err := LoadTemplate(sub)
		if err != nil {
			return nil, err
		}

		env.Templates[t.CanonicalName] = t
		env.CanonicalNames[templateAlias] = t.CanonicalName
	}

	return &env, nil
}
