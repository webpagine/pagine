package structure

import (
	"github.com/webpagine/go-pagine/util"
	"github.com/webpagine/go-pagine/vfs"
	"strings"
)

type EnvManifest struct {
	Use map[string]string `toml:"use"`
}

type Env struct {
	Root vfs.DirFS

	Templates map[string]*Template
}

func LoadEnv(root vfs.DirFS) (*Env, error) {
	var env = Env{Root: root, Templates: map[string]*Template{}}
	var manifest EnvManifest

	err := util.UnmarshalTOMLFile(root, "/env.toml", &manifest)
	if err != nil {
		return nil, err
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
