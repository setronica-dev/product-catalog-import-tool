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
	Name             string `yaml:"name"`
	HeaderRowsToSkip int    `yaml:"header_rows_to_skip"`
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
		Name:             c.Name,
		HeaderRowsToSkip: c.HeaderRowsToSkip,
	}
}
