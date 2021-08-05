package configModels

type RawCommonCatalogConfig struct {
	SourcePath string         `yaml:"source"`
	SentPath   string         `yaml:"sent"`
	Sheet      RawSheetConfig `yaml:"sheet"`
}

type RawSheetConfig struct {
	Products   string               `yaml:"products"`
	Offers     string               `yaml:"offers"`
	OfferItems RawSheetParamsConfig `yaml:"offer_items"`
	Failures   string               `yaml:"failures"`
}

type RawSheetParamsConfig struct {
	Name               string `yaml:"name"`
	HeaderColumnsCount int    `yaml:"header_columns_count"`
}

func (c *RawCommonCatalogConfig) ToConfig() *CommonConfig {
	return &CommonConfig{
		SourcePath: c.SourcePath,
		SentPath:   c.SentPath,
		Sheet: &SheetConfig{
			Products:   c.Sheet.Products,
			Offers:     c.Sheet.Offers,
			Failures:   c.Sheet.Failures,
			OfferItems: c.Sheet.OfferItems.ToConfig(),
		},
	}
}

func (c *RawSheetParamsConfig) ToConfig() *SheetParamsConfig {
	return &SheetParamsConfig{
		Name:            c.Name,
		HeaderRowsCount: c.HeaderColumnsCount,
	}
}
