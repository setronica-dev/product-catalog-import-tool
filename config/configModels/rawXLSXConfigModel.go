package configModels

type RawXLSXConfig struct {
	SourcePath string         `yaml:"source"`
	SentPath   string         `yaml:"sent"`
	Sheet      RawSheetConfig `yaml:"sheet"`
}

type RawSheetConfig struct {
	Products   RawSheetParamsConfig `yaml:"products"`
	Offers     RawSheetParamsConfig `yaml:"offers"`
	OfferItems RawSheetParamsConfig `yaml:"offer_items"`
	Attributes RawSheetParamsConfig `yaml:"attributes"`
}

type RawSheetParamsConfig struct {
	Name             string `yaml:"name"`
	HeaderRowsToSkip int    `yaml:"header_rows_to_skip"`
}

func (c *RawXLSXConfig) ToConfig() *XLSXConfig {
	if *c == (RawXLSXConfig{}) {
		return nil
	}
	return &XLSXConfig{
		SourcePath: c.SourcePath,
		SentPath:   c.SentPath,
		Sheets: &SheetsConfig{
			Products:   c.Sheet.Products.ToConfig(),
			Offers:     c.Sheet.Offers.ToConfig(),
			Attributes: c.Sheet.Attributes.ToConfig(),
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
