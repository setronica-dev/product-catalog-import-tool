package product

import (
	"fmt"
	"ts/utils"
)

type ProductColumnMap struct {
	ProductID string
	Category  string
	Name      string
}

func (p *ProductHandler) GetCurrentHeader(row map[string]interface{}) *ProductColumnMap {
	var res ProductColumnMap

	columnMap := p.ColumnMap
	for k, _ := range row {
		key := fmt.Sprintf("%v", k)
		switch utils.TrimAll(key) {
		case utils.TrimAll(columnMap.Category):
			res.Category = key
		case utils.TrimAll(columnMap.ProductID):
			res.ProductID = key
		case utils.TrimAll(columnMap.Name):
			res.Name = key
		}
	}
	return &res
}
