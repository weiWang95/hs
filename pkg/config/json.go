package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var basePath, _ = os.Getwd()

func LoadJsonConfig(file string, v any) error {
	bs, err := os.ReadFile(filepath.Join(basePath, file))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, v); err != nil {
		return err
	}

	return nil
}
