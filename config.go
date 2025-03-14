package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type Config struct {
	GitSearchDirs []string
	ExtraDirs     []string
	// TODO ignore patterns
}

func DefaultConfigFile() string {
	return path.Join(os.Getenv("HOME"), ".lazytsm.toml")
}

func ReadConfig(path string) (Config, error) {
	// Check that config file exists
	if _, err := os.Stat(path); err != nil {
		return Config{}, fmt.Errorf("error reading config at %v: %v", path, err)
	}

	var conf Config

	_, err := toml.DecodeFile(path, &conf)

	if err != nil {
		return Config{}, fmt.Errorf("Error parsing config file at %v: %v", path, err)
	}

	return conf, nil
}

func ReadDefaultConfig() (Config, error) {
	return ReadConfig(DefaultConfigFile())
}
