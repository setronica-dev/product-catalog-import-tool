package attribute

import (
	"go.uber.org/dig"
	"ts/adapters"
	"ts/productImport/reports"
)

type AttributeHandlerInterface interface {
	InitAttributeData(filePath string) ([]*Attribute, error)
}

type Deps struct {
	dig.In
	FileManager *adapters.FileManager
	Handler     adapters.HandlerInterface
	Report      *reports.ReportsHandler
}