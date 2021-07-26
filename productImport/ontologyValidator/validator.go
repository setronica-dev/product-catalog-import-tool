package ontologyValidator

import (
	"ts/productImport/attribute"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/product"
	"ts/productImport/reports"
)

type Validator struct {
	productHandler product.ProductHandlerInterface
	ColumnMap      *ColumnMap
}

type ColumnMap struct {
	Category  string
	ProductID string
	Name      string
}

func NewValidator(deps Deps) ValidatorInterface {
	m := deps.Mapper.Parse()
	return &Validator{
		productHandler: deps.ProductHandler,
		ColumnMap: &ColumnMap{
			Category:  m.Category,
			ProductID: m.ProductID,
			Name:      m.Name,
		},
	}
}

func (v *Validator) Validate(data struct {
	Mapping       map[string]string
	Rules         *models.OntologyConfig
	SourceData    []map[string]interface{}
	AttributeData []*attribute.Attribute
}) ([]reports.Report, bool) {

	parsedProducts := v.productHandler.InitParsedSourceData(data.SourceData)

	if data.AttributeData != nil && len(data.AttributeData) > 0 {
		report, isErr := v.validateReport(data.Rules, parsedProducts, data.AttributeData)
		return report, isErr
	}

	report, isErr := v.validateSource(data)
	return report, isErr
}
