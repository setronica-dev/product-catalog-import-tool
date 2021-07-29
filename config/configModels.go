package config

type ServiceConfig struct {
	Port uint16
}

type TradeshiftAPIConfig struct {
	APIBaseURL     string
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
	TenantId       string
	Currency       string
	FileLocale     string
}

type ProductCatalogConfig struct {
	SourcePath                 string
	ReportPath                 string
	SecondValidationSourcePath string
	MappingPath                string
	OntologyPath               string
	SentPath                   string
	InProgressPath             string
	SuccessResultPath          string
	FailResultPath             string
}

type OfferCatalogConfig struct {
	SourcePath string
	SentPath   string
}

type OfferItemCatalogConfig struct {
	SourcePath        string
	SuccessResultPath string
	ReportPath        string
	SentPath          string
}

type CommonConfig struct {
	SourcePath string
	SentPath   string
	Sheet      *SheetConfig
}

type SheetConfig struct {
	Products   string
	Offers     string
	Failures   string
	OfferItems *SheetParamsConfig
}

type SheetParamsConfig struct {
	Name            string
	HeaderRowsCount int
}
