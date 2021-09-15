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
	Mapper         mapping.MappingHandlerInterface
}

type ValidatorInterface interface {
	InitialValidation(mapping map[string]string,
		rules *models.OntologyConfig,
		sourceData []map[string]interface{}) ([]reports.Report, bool)
	SecondaryValidation(
		mapping map[string]string,
		rules *models.OntologyConfig,
		sourceData []map[string]interface{},
		attributeData []*attribute.Attribute) ([]reports.Report, bool)
}
