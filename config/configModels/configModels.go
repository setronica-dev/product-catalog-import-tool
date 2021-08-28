package configModels

type ServiceConfig struct {
	Port uint16
}

type ProductCatalogConfig struct {
	SourcePath        string
	ReportPath        string
	MappingPath       string
	OntologyPath      string
	SentPath          string
	InProgressPath    string
	SuccessResultPath string
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
	Sheets     *SheetsConfig
}

type SheetsConfig struct {
	Products   *SheetParamsConfig
	Offers     *SheetParamsConfig
	Failures   *SheetParamsConfig
	OfferItems *SheetParamsConfig
}

type SheetParamsConfig struct {
	Name             string
	HeaderRowsToSkip int
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
	Recipients     *Recipients
}
