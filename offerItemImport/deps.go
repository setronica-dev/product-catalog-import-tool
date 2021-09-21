package offerItemImport

import (
	"go.uber.org/dig"
	"ts/config"
	"ts/offerItemImport/offerItemMapping"
	"ts/outwardImport"
	"ts/productImport/mapping"
)

type Deps struct {
	dig.In
	OutwardImportHandler outwardImport.OutwardImportInterface
	OfferItemMapping     offerItemMapping.OfferItemMappingHandlerInterface
	Mapping              mapping.MappingHandlerInterface
	Config               *config.Config
}
