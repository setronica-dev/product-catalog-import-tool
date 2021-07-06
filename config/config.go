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

	if err := GetValidator().Struct(rawServiceConfig); err != nil {
		log.Fatalf("invalid config file %s\n%s", serviceConfigPath, err)
	}

	// unmarshal into the tmp raw config
	config = configFromRaw(rawServiceConfig)
}

type Config struct {
	Service        ServiceConfig
	ProductCatalog ProductCatalogConfig
	OfferCatalog   OfferCatalogConfig
	CommonConfig   CommonConfig
	TradeshiftAPI  TradeshiftAPIConfig
}

func Get() *Config {
	result := &Config{}
	*result = *config
	return result
}

func configFromRaw(rawService *RawServiceConfig) *Config {
	p := rawService.ProductCatalogConfig
	t := rawService.TradeshiftAPIConfig
	o := rawService.OfferCatalogConfig
	c := rawService.CommonConfig
	return &Config{
		Service:        *rawService.ToConfig(),
		ProductCatalog: *p.ToConfig(),
		OfferCatalog:   *o.ToConfig(),
		CommonConfig:   *c.ToConfig(),
		TradeshiftAPI:  *t.ToConfig(),
	}
}
