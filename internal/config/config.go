package config

import (
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
	"io/fs"
	"os"
	"path"
)

type Config struct {
	APIKey string `setting:"API Key"`
}

func ReadConfig() (Config, error) {
	c := Config{}
	cfgDir := path.Join(xdg.ConfigHome, "soramail")
	cfg := path.Join(cfgDir, "config.toml")

	_ = os.Mkdir(xdg.ConfigHome, 0755)
	err := os.Mkdir(cfgDir, 0755)

	if !errors.Is(err, fs.ErrExist) {
		return c, err
	}

	if _, err := os.Stat(cfg); errors.Is(err, fs.ErrNotExist) {
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

func SaveConfig(config Config) error {
	cfgDir := path.Join(xdg.ConfigHome, "soramail")
	cfg := path.Join(cfgDir, "config.toml")

	_ = os.Mkdir(xdg.ConfigHome, 0755)
	err := os.Mkdir(cfgDir, 0755)

	if !errors.Is(err, fs.ErrExist) {
		return err
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed encoding Config: %v", err)
	}

	if err := os.WriteFile(cfg, data, 0640); err != nil {
		return fmt.Errorf("failed creating Config: %v", err)
	}

	return nil
}
