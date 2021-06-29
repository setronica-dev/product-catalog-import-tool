package rawOntology

import (
	"ts/productImport/ontologyRead/models"
)

type RawOntology struct {
	categories map[string]*RawCategory
}

func NewRawOntology() *RawOntology {
	return &RawOntology{
		categories: map[string]*RawCategory{},
	}
}

func (o *RawOntology) AddCategoryAttribute(c *RawCategory, a *RawAttribute) error {
	_, ok := o.categories[c.unspsc]
	if !ok {
		o.categories[c.unspsc] = c
	}
	err := o.categories[c.unspsc].addAttribute(a)
	return err
}

func (o *RawOntology) ToConfig() *models.OntologyConfig {
	configs := make(map[string]*models.CategoryConfig, len(o.categories))
	for i, v := range o.categories {
		configs[i] = v.ToConfig()
	}
	return &models.OntologyConfig{
		Categories: configs,
	}
}

func (o *RawOntology) GetCategoriesCount() int {
	return len(o.categories)
}
