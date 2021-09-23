package configModels

type RawProductCatalogConfig struct {
	SourcePath                 string `yaml:"source"`
	MappingPath                string `yaml:"mapping"`
	OntologyPath               string `yaml:"ontology"`
	SentPath                   string `yaml:"sent"`
	InProgressPath             string `yaml:"in_progress"`
	ReportPath                 string `yaml:"report"`
	SuccessResultPath          string `yaml:"success_result"`
}

type RawOfferCatalogConfig struct {
	SourcePath string `yaml:"source"`
	SentPath   string `yaml:"sent"`
}

type RawOfferItemCatalogConfig struct {
	SourcePath        string `yaml:"source"`
	SuccessResultPath string `yaml:"success_result"`
	ReportPath        string `yaml:"report"`
	SentPath          string `yaml:"sent"`
}

func (c *RawProductCatalogConfig) ToConfig() *ProductCatalogConfig {
	return &ProductCatalogConfig{
		SourcePath:                 c.SourcePath,
		ReportPath:                 c.ReportPath,
		MappingPath:                c.MappingPath,
		OntologyPath:               c.OntologyPath,
		SentPath:                   c.SentPath,
		InProgressPath:             c.InProgressPath,
		SuccessResultPath:          c.SuccessResultPath,
	}
}

func (c *RawOfferCatalogConfig) ToConfig() *OfferCatalogConfig {
	return &OfferCatalogConfig{
		SourcePath: c.SourcePath,
		SentPath:   c.SentPath,
	}
}

func (c *RawOfferItemCatalogConfig) ToConfig() *OfferItemCatalogConfig {
	return &OfferItemCatalogConfig{
		SourcePath:        c.SourcePath,
		SuccessResultPath: c.SuccessResultPath,
		ReportPath:        c.ReportPath,
		SentPath:          c.SentPath,
	}
}
