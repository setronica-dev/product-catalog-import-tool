package ontologyValidator

import (
	"ts/models"
)

type ValidatorInterface interface {
	Validate(data struct {
		Mapping map[string]string
		Rules   *models.OntologyConfig
		Data    []map[string]interface{}
		Report  []*models.Report
	}) ([]models.Report, bool)
}
