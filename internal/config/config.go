package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

func (c Config) Validate() error {
	if c.Port == 0 {
		return fmt.Errorf("invalid port")
	}

	return nil
}

func Read() (Config, error) {
	cfgFile, err := os.Open("config.json")
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file: %s", err)
	}
	defer cfgFile.Close()

	var cfg Config
	js := json.NewDecoder(cfgFile)
	if err := js.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not unmarshal json: %s", err)
	}

	return cfg, nil
}
