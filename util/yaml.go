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
