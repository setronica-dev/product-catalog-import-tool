package ontologyValidator

import (
	"ts/productImport/ontologyRead/models"
	"ts/productImport/reports"
)

type ValidatorInterface interface {
	Validate(data struct {
		Mapping map[string]string
		Rules   *models.OntologyConfig
		Data    []map[string]interface{}
		Report  []*reports.Report
	}) ([]reports.Report, bool)
}
