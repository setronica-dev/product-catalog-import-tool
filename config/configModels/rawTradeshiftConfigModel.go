package configModels

type RawTradeshiftAPIConfig struct {
	APIBaseURL     string         `yaml:"base_url" validate:"required"`
	ConsumerKey    string         `yaml:"consumer_key" validate:"required"`
	ConsumerSecret string         `yaml:"consumer_secret" validate:"required"`
	Token          string         `yaml:"token" validate:"required"`
	TokenSecret    string         `yaml:"token_secret" validate:"required"`
	TenantId       string         `yaml:"tenant_id" validate:"required"`
	Currency       string         `yaml:"currency" validate:"required"`
	FileLocale     string         `yaml:"file_locale" validate:"required"`
	Recipients     []RawRecipient `yaml:"recipients"`
}

type RawRecipient struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

func (r *RawTradeshiftAPIConfig) ToConfig() *TradeshiftAPIConfig {
	recipients := make([]*recipient, 0)
	for _, item := range r.Recipients {
		recipients = append(recipients, item.ToConfig())
	}

	return &TradeshiftAPIConfig{
		APIBaseURL:     r.APIBaseURL,
		ConsumerKey:    r.ConsumerKey,
		ConsumerSecret: r.ConsumerSecret,
		Token:          r.Token,
		TokenSecret:    r.TokenSecret,
		TenantId:       r.TenantId,
		Currency:       r.Currency,
		FileLocale:     r.FileLocale,
		Recipients: &Recipients{
			collection: recipients,
		},
	}
}

func (r *RawRecipient) ToConfig() *recipient {
	return &recipient{
		id:   r.ID,
		name: r.Name,
	}
}
