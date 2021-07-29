package importToTradeshift

import (
	"go.uber.org/dig"
	"ts/adapters"
	"ts/config"
	"ts/externalAPI/tradeshiftAPI"
	outwardImport "ts/outwardImport"
)

type TradeshiftImport struct {
	transport   *tradeshiftAPI.TradeshiftAPI
	filemanager *adapters.FileManager
	handler     adapters.HandlerInterface
	tsConfig    *TradeShiftConfiguration
}

type Deps struct {
	dig.In
	Config        *config.Config
	TradeshiftAPI *tradeshiftAPI.TradeshiftAPI
	FileManager   *adapters.FileManager
	FilesHandler  adapters.HandlerInterface
}

type TradeShiftConfiguration struct {
	tsCurrency string
	tsLocale   string
}

func NewImportToTradeshift(deps Deps) outwardImport.OutwardImportInterface{
	h := deps.FilesHandler
	h.Init(adapters.TXT)

	return &TradeshiftImport{
		transport:   deps.TradeshiftAPI,
		filemanager: deps.FileManager,
		handler:     h,
		tsConfig: &TradeShiftConfiguration{
			tsCurrency: deps.Config.TradeshiftAPI.Currency,
			tsLocale:   deps.Config.TradeshiftAPI.FileLocale,
		},
	}
}
