package configModels

type RawServiceConfig struct {
	Port                   uint16                    `yaml:"port"`
	ProductCatalogConfig   RawProductCatalogConfig   `yaml:"product" validate:"required"`
	OfferCatalogConfig     RawOfferCatalogConfig     `yaml:"offer" validate:"required"`
	OfferItemCatalogConfig RawOfferItemCatalogConfig `yaml:"offer_item" validate:"required"`
	TradeshiftAPIConfig    RawTradeshiftAPIConfig    `yaml:"tradeshift_api" validate:"required"`
	XLSXConfig             RawXLSXConfig             `yaml:"xlsx_settings"`
}

func (c *RawServiceConfig) ToConfig() *ServiceConfig {
	return &ServiceConfig{
		Port: c.Port,
	}
}
