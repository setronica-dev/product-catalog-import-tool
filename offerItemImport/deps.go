package offerItemImport

import (
	"go.uber.org/dig"
	"ts/config"
	"ts/outwardImport"
	"ts/productImport/mapping"
)

type Deps struct {
	dig.In
	OutwardImportHandler outwardImport.OutwardImportInterface
	Mapping              mapping.MappingHandlerInterface
	Config               *config.Config
}
