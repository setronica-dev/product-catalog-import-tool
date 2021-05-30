package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const (
	serviceConfigPath = "./service.yaml"
)

var config *Config

func Init() {
	// read config data
	data, err := ioutil.ReadFile(serviceConfigPath)
	if err != nil {
		log.Fatalf("unable to load config file from %s\n%s", serviceConfigPath, err)
	}

	// unmarshal into the tmp raw config
	rawServiceConfig := &RawServiceConfig{}
	if err := yaml.Unmarshal(data, rawServiceConfig); err != nil {
		log.Fatalf("unable to unmarshal config file %s\n%s", serviceConfigPath, err)
	}

	// unmarshal into the tmp raw config
	config = configFromRaw(rawServiceConfig)
}

type Config struct {
	Service       ServiceConfig
	Catalog       CatalogConfig
	TradeshiftAPI TradeshiftAPIConfig
}

func Get() *Config {
	result := &Config{}
	*result = *config
	return result
}

func configFromRaw(rawService *RawServiceConfig) *Config {
	c := rawService.CatalogConfig
	t := rawService.TradeshiftAPIConfig
	return &Config{
		Service:       *rawService.ToConfig(),
		Catalog:       *c.ToConfig(),
		TradeshiftAPI: *t.ToConfig(),
	}
}
