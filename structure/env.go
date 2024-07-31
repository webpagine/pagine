// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/webpagine/pagine/v2/util"
	"github.com/webpagine/pagine/v2/vfs"
	"regexp"
	"strings"
)

type EnvManifest struct {
	Use map[string]string `yaml:"use"`

	Ignore []string `yaml:"ignore"`
}

type Env struct {
	Root *vfs.DirFS

	Templates map[string]*Template

	CanonicalNames map[string]string

	IgnoreGlobs []*regexp.Regexp
}

func (e *Env) BaseOf(fs *vfs.DirFS) string {
	base, _ := strings.CutPrefix(fs.Path, e.Root.Path)
	return base
}

func LoadEnv(root *vfs.DirFS) (*Env, error) {
	var env = Env{Root: root, Templates: map[string]*Template{}}
	var manifest EnvManifest

	err := util.UnmarshalYAMLFile(root, "/env.yaml", &manifest)
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

	for templateName, templatePath := range manifest.Use {
		sub, err := root.Chroot(templatePath)
		if err != nil {
			return nil, err
		}

		t, err := LoadTemplate(sub)
		if err != nil {
			return nil, err
		}

		env.Templates[templateName] = t
		env.CanonicalNames[t.CanonicalName] = templateName
	}

	return &env, nil
}

func ParseTemplatePair(pair string) (string, string) {
	split := strings.Split(pair, ":")
	if len(split) == 2 {
		return split[0], split[1]
	}

	return pair, "main"
}
