package mapping

import (
	"go.uber.org/dig"
	"ts/config"
)

type ColumnMap struct {
	Category  string
	ProductID string
	Name      string
}

type MappingHandlerInterface interface {
	Init(mappingConfigPath string) map[string]string
	Get() map[string]string
	Parse() *ColumnMap
}

type Deps struct {
	dig.In
	Config *config.Config
}
