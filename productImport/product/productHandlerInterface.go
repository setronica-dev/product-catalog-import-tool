package product

import (
	"go.uber.org/dig"
	"ts/adapters"
	"ts/productImport/mapping"
)

type ProductHandlerInterface interface {
	InitSourceData(sourceFeedPath string) ([]map[string]interface{}, error)
	GetCurrentHeader(row map[string]interface{}) *ColumnMap
}

type Deps struct {
	dig.In
	FileManager *adapters.FileManager
	Handler     adapters.HandlerInterface
	Mapping     mapping.MappingHandlerInterface
}
