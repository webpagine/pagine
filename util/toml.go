package util

import (
	"github.com/BurntSushi/toml"
	"os"
)

func UnmarshalTOMLFile(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return toml.Unmarshal(b, v)
}
