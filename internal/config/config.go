package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Broker struct {
		Handlers   int `json:"event_handlers"`
		Buffer     int `json:"event_buffer"`
		MaxRetries int `json:"max_retries"`
	} `json:"broker"`
	InvoiceDB  string `json:"invoice_db"`
	IssuerDB   string `json:"issuer_db"`
	InvestorDB string `json:"investor_db"`
}

func (c Config) Validate() error {
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
