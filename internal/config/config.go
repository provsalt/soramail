package config

import (
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
	"os"
	"path"
)

type Config struct {
	APIKey string
}

func ReadConfig() (Config, error) {
	c := Config{}
	cfgDir := path.Join(xdg.ConfigHome, "soramail")
	cfg := path.Join(cfgDir, "config.toml")

	err := os.Mkdir(cfgDir, 0750)

	if !errors.Is(err, os.ErrExist) {
		return c, err
	}

	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default Config: %v", err)
		}
		if err := os.WriteFile(cfg, data, 0640); err != nil {
			return c, fmt.Errorf("failed creating Config: %v", err)
		}
		return c, nil
	}
	data, err := os.ReadFile(cfg)
	if err != nil {
		return c, fmt.Errorf("error reading Config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding Config: %v", err)
	}
	return c, nil

}
