package config

type ServiceConfig struct {
	Port                       uint16
}

type TradeshiftAPIConfig struct {
	APIBaseURL     string
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
	TenantId       string
}

type CatalogConfig struct {
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
