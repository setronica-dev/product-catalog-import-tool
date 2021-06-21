package rawOntology

import (
	"fmt"
	"ts/productImport/ontologyRead/models"
)

type RawCategory struct {
	unspsc     string
	unspscName string
	attributes map[string]*RawAttribute
}

func NewRawCategory(raw map[string]interface{}, header *RawHeader) *RawCategory {
	c := RawCategory{
		unspsc:     raw[header.category].(string),
		attributes: map[string]*RawAttribute{},
	}
	if header.categoryName != "" {
		c.unspscName = raw[header.categoryName].(string)
	}
	return &c
}

func (c *RawCategory) addAttribute(a *RawAttribute) error {
	_, ok := c.attributes[a.name]
	if !ok {
		c.attributes[a.name] = a
		return nil
	} else {
		return fmt.Errorf("attribute %v is allready exists in category %v", a.name, c.unspsc)
	}
}

func (c *RawCategory) ToConfig() *models.CategoryConfig {
	configs := make(map[string]*models.AttributeConfig, len(c.attributes))
	for i, v := range c.attributes {
		configs[i] = v.ToConfig()
	}
	return &models.CategoryConfig{
		UNSPSC:     c.unspsc,
		Name:       c.unspscName,
		Attributes: configs,
	}
}
