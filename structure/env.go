package structure

import (
	"github.com/webpagine/pagine/util"
	"github.com/webpagine/pagine/vfs"
	"regexp"
	"strings"
)

type EnvManifest struct {
	Use map[string]string `toml:"use"`

	Ignore []string `toml:"ignore"`
}

type Env struct {
	Root vfs.DirFS

	Templates map[string]*Template

	IgnoreGlobs []*regexp.Regexp
}

func LoadEnv(root vfs.DirFS) (*Env, error) {
	var env = Env{Root: root, Templates: map[string]*Template{}}
	var manifest EnvManifest

	err := util.UnmarshalTOMLFile(root, "/env.toml", &manifest)
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
		t, err := LoadTemplate(root.DirFS(templatePath))
		if err != nil {
			return nil, err
		}

		env.Templates[templateName] = t
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
