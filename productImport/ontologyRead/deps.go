package ontologyRead

import (
	"go.uber.org/dig"
	"ts/adapters"
	"ts/config"
)

type Deps struct {
	dig.In
	Config       *config.Config
	Handler      adapters.HandlerInterface
	FilesManager *adapters.FileManager
}
