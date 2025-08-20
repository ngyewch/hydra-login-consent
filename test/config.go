package main

import (
	"path/filepath"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func mergeConfig(k *koanf.Koanf, configFile string) error {
	if configFile == "" {
		return nil
	}
	ext := filepath.Ext(configFile)
	var parser koanf.Parser
	switch ext {
	case ".json":
		parser = json.Parser()
	case ".toml":
		parser = toml.Parser()
	case ".yml":
		parser = yaml.Parser()
	case ".yaml":
		parser = yaml.Parser()
	}
	err := k.Load(file.Provider(configFile), parser)
	if err != nil {
		return err
	}
	return nil
}
