package offerItemMapping

import (
	"go.uber.org/dig"
	"ts/config"
	"ts/outwardImport"
	"ts/productImport/mapping"
)

type OfferItemMappingHandlerInterface interface {
	Run() error
}

type Deps struct {
	dig.In
	OutwardImportHandler outwardImport.OutwardImportInterface
	Mapping              mapping.MappingHandlerInterface
	Config               *config.Config
}
