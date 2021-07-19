package ontologyValidator

import (
	"go.uber.org/dig"
	"ts/productImport/attribute"
	"ts/productImport/mapping"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/product"
	"ts/productImport/reports"
)

type Deps struct {
	dig.In
	ProductHandler product.ProductHandlerInterface
	Mapper mapping.MappingHandlerInterface
}

type ValidatorInterface interface {
	Validate(data struct {
		Mapping       map[string]string
		Rules         *models.OntologyConfig
		SourceData    []map[string]interface{}
		AttributeData []*attribute.Attribute
	}) ([]reports.Report, bool)
}
