package configModels

type RawCommonCatalogConfig struct {
	SourcePath string         `yaml:"source"`
	SentPath   string         `yaml:"sent"`
	Sheet      RawSheetConfig `yaml:"sheet"`
}

type RawSheetConfig struct {
	Products   RawSheetParamsConfig `yaml:"products"`
	Offers     RawSheetParamsConfig `yaml:"offers"`
	OfferItems RawSheetParamsConfig `yaml:"offer_items"`
	Failures   RawSheetParamsConfig `yaml:"failures"`
}

type RawSheetParamsConfig struct {
	Name             string `yaml:"name"`
	HeaderRowsToSkip int    `yaml:"header_rows_to_skip"`
}

func (c *RawCommonCatalogConfig) ToConfig() *CommonConfig {
	return &CommonConfig{
		SourcePath: c.SourcePath,
		SentPath:   c.SentPath,
		Sheets: &SheetsConfig{
			Products:   c.Sheet.Products.ToConfig(),
			Offers:     c.Sheet.Offers.ToConfig(),
			Failures:   c.Sheet.Failures.ToConfig(),
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
