package config

type RawServiceConfig struct {
	Port                 uint16                  `yaml:"port" validate:"required"`
	ProductCatalogConfig RawProductCatalogConfig `yaml:"product" validate:"required"`
	OfferCatalogConfig   RawOfferCatalogConfig   `yaml:"offer" validate:"required"`
	TradeshiftAPIConfig  RawTradeshiftAPIConfig  `yaml:"tradeshift_api" validate:"required"`
}

type RawProductCatalogConfig struct {
	SourcePath                 string `yaml:"source"`
	ReportPath                 string `yaml:"report"`
	SecondValidationSourcePath string `yaml:"source2"`
	MappingPath                string `yaml:"mapping"`
	OntologyPath               string `yaml:"ontology"`
	SentPath                   string `yaml:"sent"`
	InProgressPath             string `yaml:"in_progress"`
	SuccessResultPath          string `yaml:"success_result"`
	FailResultPath             string `yaml:"fail_result"`
}

type RawOfferCatalogConfig struct {
	SourcePath string `yaml:"source"`
	SentPath   string `yaml:"sent"`
}

type RawTradeshiftAPIConfig struct {
	APIBaseURL     string `yaml:"base_url" validate:"required"`
	ConsumerKey    string `yaml:"consumer_key" validate:"required"`
	ConsumerSecret string `yaml:"consumer_secret" validate:"required"`
	Token          string `yaml:"token" validate:"required"`
	TokenSecret    string `yaml:"token_secret" validate:"required"`
	TenantId       string `yaml:"tenant_id" validate:"required"`
}

func (c *RawServiceConfig) ToConfig() *ServiceConfig {
	return &ServiceConfig{
		Port: c.Port,
	}
}

func (c *RawProductCatalogConfig) ToConfig() *ProductCatalogConfig {
	return &ProductCatalogConfig{
		SourcePath:                 c.SourcePath,
		ReportPath:                 c.ReportPath,
		SecondValidationSourcePath: c.SecondValidationSourcePath,
		MappingPath:                c.MappingPath,
		OntologyPath:               c.OntologyPath,
		SentPath:                   c.SentPath,
		InProgressPath:             c.InProgressPath,
		SuccessResultPath:          c.SuccessResultPath,
		FailResultPath:             c.FailResultPath,
	}
}

func (c *RawOfferCatalogConfig) ToConfig() *OfferCatalogConfig {
	return &OfferCatalogConfig{
		SourcePath: c.SourcePath,
		SentPath:   c.SentPath,
	}
}

func (t *RawTradeshiftAPIConfig) ToConfig() *TradeshiftAPIConfig {
	return &TradeshiftAPIConfig{
		APIBaseURL:     t.APIBaseURL,
		ConsumerKey:    t.ConsumerKey,
		ConsumerSecret: t.ConsumerSecret,
		Token:          t.Token,
		TokenSecret:    t.TokenSecret,
		TenantId:       t.TenantId,
	}
}
