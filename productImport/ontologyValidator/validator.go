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
	m := deps.Mapper.GetColumnMapConfig()
	return &Validator{
		productHandler: deps.ProductHandler,
		ColumnMap: &ColumnMap{
			Category:  m.Category,
			ProductID: m.ProductID,
			Name:      m.Name,
		},
	}
}

func (v *Validator) InitialValidation(
	mapping map[string]string,
	rules *models.OntologyConfig,
	sourceData []map[string]interface{}) ([]reports.Report, bool) {
	report, isErr := v.validateProductsAgainstRules(mapping,
		rules,
		sourceData,
	)
	return report, isErr
}

func (v *Validator) SecondaryValidation(
	mapping map[string]string,
	rules *models.OntologyConfig,
	sourceData []map[string]interface{},
	attributeData []*attribute.Attribute,
) ([]reports.Report, bool) {

	parsedProducts := v.productHandler.InitParsedSourceData(sourceData)
	if attributeData != nil && len(attributeData) > 0 {
		report, isErr := v.validateAttributesAgainstRules(rules, parsedProducts, attributeData)
		return report, isErr
	}

	report, isErr := v.validateProductsAgainstRules(mapping,
		rules,
		sourceData,
	)
	return report, isErr
}
