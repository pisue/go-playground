package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
}

func NewConfig(path string) *Config {
	c := new(Config)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = toml.NewDecoder(file).Decode(c); err != nil {
		panic(err)
	}

	return c
}
